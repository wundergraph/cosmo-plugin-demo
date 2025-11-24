import * as grpc from '@grpc/grpc-js';
import { PluginServer } from './plugin-server';

// Import generated gRPC code
import { 
  CoursesServiceService, 
  ICoursesServiceServer 
} from '../generated/service_grpc_pb.js';
import { 
  LookupUserByIdRequest,
  LookupUserByIdResponse,
  QueryCoursesRequest,
  QueryCoursesResponse,
  QueryCourseRequest,
  QueryCourseResponse,
  MutationCreateCourseRequest,
  MutationCreateCourseResponse,
  MutationEnrollUserRequest,
  MutationEnrollUserResponse,
  User,
  Course,
  Enrollment
} from '../generated/service_pb.js';

// Thread-safe counters for generating unique IDs
const courseCounterBuffer = new SharedArrayBuffer(4);
const courseCounterArray = new Int32Array(courseCounterBuffer);
Atomics.store(courseCounterArray, 0, 0);

const enrollmentCounterBuffer = new SharedArrayBuffer(4);
const enrollmentCounterArray = new Int32Array(enrollmentCounterBuffer);
Atomics.store(enrollmentCounterArray, 0, 0);

// In-memory data stores
const courses = new Map<string, Course>();
const enrollments = new Map<string, Enrollment>();
const userEnrollments = new Map<string, string[]>(); // userId -> enrollmentIds
const courseEnrollments = new Map<string, string[]>(); // courseId -> enrollmentIds
const userCourses = new Map<string, string[]>(); // instructorId -> courseIds

// Seed initial mock data
function seedMockData() {
  console.error('Seeding mock course data...');

  // Create mock instructors
  const instructor1 = new User();
  instructor1.setId('1'); // Corresponds to user ID 1 from users plugin

  const instructor2 = new User();
  instructor2.setId('2'); // Corresponds to user ID 2 from users plugin

  // Create mock courses
  const course1 = new Course();
  course1.setId('course-1');
  course1.setTitle('Introduction to GraphQL Federation');
  course1.setDescription('Learn the fundamentals of GraphQL Federation and how to build distributed GraphQL architectures');
  course1.setInstructor(instructor1);
  course1.setDurationHours(20);
  course1.setPublished(true);
  course1.setEnrollmentsList([]);

  const course2 = new Course();
  course2.setId('course-2');
  course2.setTitle('Advanced TypeScript Patterns');
  course2.setDescription('Master advanced TypeScript patterns including generics, decorators, and type manipulation');
  course2.setInstructor(instructor1);
  course2.setDurationHours(15);
  course2.setPublished(true);
  course2.setEnrollmentsList([]);

  const course3 = new Course();
  course3.setId('course-3');
  course3.setTitle('Building Scalable APIs with gRPC');
  course3.setDescription('Learn how to design and implement high-performance APIs using gRPC and Protocol Buffers');
  course3.setInstructor(instructor2);
  course3.setDurationHours(25);
  course3.setPublished(true);
  course3.setEnrollmentsList([]);

  const course4 = new Course();
  course4.setId('course-4');
  course4.setTitle('React Performance Optimization');
  course4.setDescription('Deep dive into React performance optimization techniques and best practices');
  course4.setInstructor(instructor2);
  course4.setDurationHours(12);
  course4.setPublished(false);
  course4.setEnrollmentsList([]);

  // Store courses
  courses.set('course-1', course1);
  courses.set('course-2', course2);
  courses.set('course-3', course3);
  courses.set('course-4', course4);

  // Track instructor courses
  userCourses.set('1', ['course-1', 'course-2']);
  userCourses.set('2', ['course-3', 'course-4']);

  // Create mock enrollments
  const user3 = new User();
  user3.setId('3');

  const user4 = new User();
  user4.setId('4');

  const user5 = new User();
  user5.setId('5');

  // Enrollment 1: User 3 enrolled in course 1 (50% complete)
  const enrollment1 = new Enrollment();
  enrollment1.setId('enrollment-1');
  enrollment1.setUser(user3);
  enrollment1.setCourse(course1);
  enrollment1.setProgress(50);
  enrollment1.setEnrolledAt('2024-01-15T10:30:00Z');

  // Enrollment 2: User 3 enrolled in course 3 (25% complete)
  const enrollment2 = new Enrollment();
  enrollment2.setId('enrollment-2');
  enrollment2.setUser(user3);
  enrollment2.setCourse(course3);
  enrollment2.setProgress(25);
  enrollment2.setEnrolledAt('2024-02-01T14:20:00Z');

  // Enrollment 3: User 4 enrolled in course 1 (75% complete)
  const enrollment3 = new Enrollment();
  enrollment3.setId('enrollment-3');
  enrollment3.setUser(user4);
  enrollment3.setCourse(course1);
  enrollment3.setProgress(75);
  enrollment3.setEnrolledAt('2024-01-20T09:15:00Z');

  // Enrollment 4: User 4 enrolled in course 2 (100% complete)
  const enrollment4 = new Enrollment();
  enrollment4.setId('enrollment-4');
  enrollment4.setUser(user4);
  enrollment4.setCourse(course2);
  enrollment4.setProgress(100);
  enrollment4.setEnrolledAt('2024-01-10T11:00:00Z');

  // Enrollment 5: User 5 enrolled in course 3 (10% complete)
  const enrollment5 = new Enrollment();
  enrollment5.setId('enrollment-5');
  enrollment5.setUser(user5);
  enrollment5.setCourse(course3);
  enrollment5.setProgress(10);
  enrollment5.setEnrolledAt('2024-03-05T16:45:00Z');

  // Store enrollments
  enrollments.set('enrollment-1', enrollment1);
  enrollments.set('enrollment-2', enrollment2);
  enrollments.set('enrollment-3', enrollment3);
  enrollments.set('enrollment-4', enrollment4);
  enrollments.set('enrollment-5', enrollment5);

  // Track user enrollments
  userEnrollments.set('3', ['enrollment-1', 'enrollment-2']);
  userEnrollments.set('4', ['enrollment-3', 'enrollment-4']);
  userEnrollments.set('5', ['enrollment-5']);

  // Track course enrollments
  courseEnrollments.set('course-1', ['enrollment-1', 'enrollment-3']);
  courseEnrollments.set('course-2', ['enrollment-4']);
  courseEnrollments.set('course-3', ['enrollment-2', 'enrollment-5']);

  // Update counters to avoid ID conflicts
  Atomics.store(courseCounterArray, 0, 4);
  Atomics.store(enrollmentCounterArray, 0, 5);

  console.error('Mock data seeded: 4 courses, 5 enrollments');
}

// Define the service implementation using the generated types
const CoursesServiceImplementation: ICoursesServiceServer = {
  // Entity resolution for User type
  lookupUserById: (
    call: grpc.ServerUnaryCall<LookupUserByIdRequest, LookupUserByIdResponse>, 
    callback: grpc.sendUnaryData<LookupUserByIdResponse>
  ) => {
    const keys = call.request.getKeysList();
    const response = new LookupUserByIdResponse();

    // For each requested user ID, return user with their courses and enrollments
    keys.forEach((key) => {
      const userId = key.getId();
      const user = new User();
      user.setId(userId);

      // Get courses where this user is instructor
      const instructorCourseIds = userCourses.get(userId) || [];
      const instructorCourses = instructorCourseIds
        .map(id => courses.get(id))
        .filter(c => c !== undefined) as Course[];
      user.setInstructorCoursesList(instructorCourses);

      // Get user's enrollments
      const enrollmentIds = userEnrollments.get(userId) || [];
      const userEnrollmentList = enrollmentIds
        .map(id => enrollments.get(id))
        .filter(e => e !== undefined) as Enrollment[];
      user.setEnrollmentsList(userEnrollmentList);

      response.addResult(user);
    });

    callback(null, response);
  },

  // Query all courses
  queryCourses: (
    call: grpc.ServerUnaryCall<QueryCoursesRequest, QueryCoursesResponse>, 
    callback: grpc.sendUnaryData<QueryCoursesResponse>
  ) => {
    const response = new QueryCoursesResponse();
    const allCourses = Array.from(courses.values());
    response.setCoursesList(allCourses);
    callback(null, response);
  },

  // Query single course by ID
  queryCourse: (
    call: grpc.ServerUnaryCall<QueryCourseRequest, QueryCourseResponse>, 
    callback: grpc.sendUnaryData<QueryCourseResponse>
  ) => {
    const courseId = call.request.getId();
    const course = courses.get(courseId);
    
    const response = new QueryCourseResponse();
    if (course) {
      response.setCourse(course);
    }
    callback(null, response);
  },

  // Create a new course
  mutationCreateCourse: (
    call: grpc.ServerUnaryCall<MutationCreateCourseRequest, MutationCreateCourseResponse>, 
    callback: grpc.sendUnaryData<MutationCreateCourseResponse>
  ) => {
    const input = call.request.getInput();
    if (!input) {
      return callback({
        code: grpc.status.INVALID_ARGUMENT,
        message: 'Input is required'
      } as grpc.ServiceError);
    }

    const courseId = `course-${Atomics.add(courseCounterArray, 0, 1) + 1}`;
    const instructorId = input.getInstructorId();

    // Create instructor user
    const instructor = new User();
    instructor.setId(instructorId);

    // Create the course
    const course = new Course();
    course.setId(courseId);
    course.setTitle(input.getTitle());
    course.setDescription(input.getDescription());
    course.setInstructor(instructor);
    course.setDurationHours(input.getDurationHours());
    course.setPublished(input.getPublished());
    course.setEnrollmentsList([]);

    // Store the course
    courses.set(courseId, course);

    // Track instructor's courses
    if (!userCourses.has(instructorId)) {
      userCourses.set(instructorId, []);
    }
    userCourses.get(instructorId)!.push(courseId);

    const response = new MutationCreateCourseResponse();
    response.setCreateCourse(course);
    callback(null, response);
  },

  // Enroll a user in a course
  mutationEnrollUser: (
    call: grpc.ServerUnaryCall<MutationEnrollUserRequest, MutationEnrollUserResponse>, 
    callback: grpc.sendUnaryData<MutationEnrollUserResponse>
  ) => {
    const userId = call.request.getUserId();
    const courseId = call.request.getCourseId();

    // Check if course exists
    const course = courses.get(courseId);
    if (!course) {
      return callback({
        code: grpc.status.NOT_FOUND,
        message: `Course ${courseId} not found`
      } as grpc.ServiceError);
    }

    // Create enrollment
    const enrollmentId = `enrollment-${Atomics.add(enrollmentCounterArray, 0, 1) + 1}`;
    
    const user = new User();
    user.setId(userId);

    const enrollment = new Enrollment();
    enrollment.setId(enrollmentId);
    enrollment.setUser(user);
    enrollment.setCourse(course);
    enrollment.setProgress(0);
    enrollment.setEnrolledAt(new Date().toISOString());

    // Store enrollment
    enrollments.set(enrollmentId, enrollment);

    // Track user enrollments
    if (!userEnrollments.has(userId)) {
      userEnrollments.set(userId, []);
    }
    userEnrollments.get(userId)!.push(enrollmentId);

    // Track course enrollments
    if (!courseEnrollments.has(courseId)) {
      courseEnrollments.set(courseId, []);
    }
    courseEnrollments.get(courseId)!.push(enrollmentId);

    const response = new MutationEnrollUserResponse();
    response.setEnrollUser(enrollment);
    callback(null, response);
  }
};

function run() {
  // Seed mock data before starting the server
  seedMockData();

  // Create the plugin server (health check automatically initialized)
  const pluginServer = new PluginServer();
  
  // Add the CoursesService service
  pluginServer.addService(CoursesServiceService, CoursesServiceImplementation);

  // Start the server
  pluginServer.serve().catch((error) => {
    console.error('Failed to start plugin server:', error);
    process.exit(1);
  });
}

run();
