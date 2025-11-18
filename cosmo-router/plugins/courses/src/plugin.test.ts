import { describe, test, expect, beforeAll, afterAll } from "bun:test";
import * as grpc from "@grpc/grpc-js";
import type { Subprocess } from "bun";

// Generated gRPC types
import { CoursesServiceClient } from '../generated/service_grpc_pb.js';
import {
  QueryCoursesRequest,
  QueryCoursesResponse,
  QueryCourseRequest,
  QueryCourseResponse,
  MutationCreateCourseRequest,
  MutationCreateCourseResponse,
  MutationEnrollUserRequest,
  MutationEnrollUserResponse,
  LookupUserByIdRequest,
  LookupUserByIdResponse,
  LookupUserByIdRequestKey,
  CourseInput
} from "../generated/service_pb.js";

describe("Query Courses", () => {
  testWithContext("should return 4 initial courses with correct structure", async (ctx) => {
    const resp = await ctx.queryCourses();
    const courses = resp.getCoursesList();
    
    expect(courses).toHaveLength(4);
    
    // Verify all courses have required fields
    courses.forEach(course => {
      expect(course.getId()).toBeTruthy();
      expect(course.getTitle()).toBeTruthy();
      expect(course.getInstructor()).toBeTruthy();
      expect(course.getDurationHours()).toBeGreaterThan(0);
    });
    
    // Spot check one course
    const firstCourse = courses.find(c => c.getId() === 'course-1');
    expect(firstCourse).toBeTruthy();
    expect(firstCourse!.getTitle()).toBe('Introduction to GraphQL Federation');
  });
});

describe("Query Single Course", () => {
  testWithContext("should return course by ID", async (ctx) => {
    const resp = await ctx.queryCourse('course-1');
    
    expect(resp.hasCourse()).toBe(true);
    const course = resp.getCourse()!;
    expect(course.getId()).toBe('course-1');
    expect(course.getTitle()).toBe('Introduction to GraphQL Federation');
    expect(course.getInstructor()!.getId()).toBe('1');
  });

  testWithContext("should return empty response for non-existent course", async (ctx) => {
    const resp = await ctx.queryCourse('non-existent-course-id');
    expect(resp.hasCourse()).toBe(false);
  });
});

describe("Query Lookup Users", () => {
  testWithContext("should return instructors with their courses", async (ctx) => {
    const resp = await ctx.lookupUserById(['1', '2']);
    const users = resp.getResultList();
    
    expect(users).toHaveLength(2);
    
    // Each instructor should have 2 courses
    users.forEach(user => {
      expect(user.getInstructorCoursesList()).toHaveLength(2);
    });
  });

  testWithContext("should return students with their enrollments", async (ctx) => {
    const resp = await ctx.lookupUserById(['3', '4']);
    const users = resp.getResultList();
    
    expect(users).toHaveLength(2);
    
    const user3 = users.find(u => u.getId() === '3')!;
    expect(user3.getEnrollmentsList()).toHaveLength(2);
    
    const user4 = users.find(u => u.getId() === '4')!;
    expect(user4.getEnrollmentsList()).toHaveLength(3);
  });

  testWithContext("should return user with no courses or enrollments", async (ctx) => {
    const resp = await ctx.lookupUserById(['999']);
    const users = resp.getResultList();
    
    expect(users).toHaveLength(1);
    expect(users[0].getInstructorCoursesList()).toHaveLength(0);
    expect(users[0].getEnrollmentsList()).toHaveLength(0);
  });
});

describe("Create Course", () => {
  testWithContext("should create course with correct fields and sequential IDs", async (ctx) => {
    const resp1 = await ctx.createCourse(
      "TypeScript Basics",
      "Learn TypeScript from scratch",
      "instructor-101",
      10,
      true
    );
    
    const course1 = resp1.getCreateCourse()!;
    expect(course1.getId()).toBe('course-5'); // First created after mock data
    expect(course1.getTitle()).toBe('TypeScript Basics');
    expect(course1.getDescription()).toBe('Learn TypeScript from scratch');
    expect(course1.getInstructor()!.getId()).toBe('instructor-101');
    expect(course1.getDurationHours()).toBe(10);
    expect(course1.getPublished()).toBe(true);
    
    // Verify sequential IDs
    const resp2 = await ctx.createCourse(
      "Advanced React",
      "Master React patterns",
      "instructor-102",
      20,
      false
    );
    const course2 = resp2.getCreateCourse()!;
    expect(course2.getId()).toBe('course-6');
    expect(course2.getPublished()).toBe(false);
  });
});


describe("Verify Created Courses Appear in Queries", () => {
  testWithContext("should find created courses in queries and lookups", async (ctx) => {
    // Verify initial count
    const initialResp = await ctx.queryCourses();
    expect(initialResp.getCoursesList()).toHaveLength(4);
    
    // Create courses
    const instructorId = 'instructor-200';
    const createResp = await ctx.createCourse(
      "New Course",
      "New Description",
      instructorId,
      15,
      true
    );
    const courseId = createResp.getCreateCourse()!.getId();
    
    // Verify in courses list
    const coursesResp = await ctx.queryCourses();
    expect(coursesResp.getCoursesList()).toHaveLength(5);
    const foundInList = coursesResp.getCoursesList().find(c => c.getId() === courseId);
    expect(foundInList).toBeTruthy();
    
    // Verify by ID query
    const byIdResp = await ctx.queryCourse(courseId);
    expect(byIdResp.hasCourse()).toBe(true);
    expect(byIdResp.getCourse()!.getTitle()).toBe('New Course');
    
    // Verify in instructor lookup
    const lookupResp = await ctx.lookupUserById([instructorId]);
    expect(lookupResp.getResultList()[0].getInstructorCoursesList()).toHaveLength(1);
    expect(lookupResp.getResultList()[0].getInstructorCoursesList()[0].getId()).toBe(courseId);
  });
});


describe("Enroll User", () => {
  testWithContext("should enroll users with sequential IDs", async (ctx) => {
    // Create a course
    const courseResp = await ctx.createCourse(
      "Popular Course",
      "Many students",
      "instructor-301",
      10,
      true
    );
    const courseId = courseResp.getCreateCourse()!.getId();
    
    // Enroll multiple users
    const resp1 = await ctx.enrollUser('user-401', courseId);
    const enroll1 = resp1.getEnrollUser()!;
    expect(enroll1.getId()).toBe('enrollment-6'); // First after mock data
    expect(enroll1.getUser()!.getId()).toBe('user-401');
    expect(enroll1.getCourse()!.getId()).toBe(courseId);
    expect(enroll1.getProgress()).toBe(0);
    
    const resp2 = await ctx.enrollUser('user-402', courseId);
    expect(resp2.getEnrollUser()!.getId()).toBe('enrollment-7');
    
    const resp3 = await ctx.enrollUser('user-403', courseId);
    expect(resp3.getEnrollUser()!.getId()).toBe('enrollment-8');
  });

  testWithContext("should fail when enrolling in non-existent course", async (ctx) => {
    await expect(
      ctx.enrollUser('user-304', 'non-existent-course-id')
    ).rejects.toThrow();
  });
});


describe("Verify Enrollments Appear in User Lookups", () => {
  testWithContext("should show enrollments in user lookups", async (ctx) => {
    const userId = 'user-501';
    
    // Create courses and enroll
    const courseIds: string[] = [];
    for (let i = 0; i < 3; i++) {
      const resp = await ctx.createCourse(
        `Course ${i}`,
        `Description ${i}`,
        `instructor-501-${i}`,
        10,
        true
      );
      courseIds.push(resp.getCreateCourse()!.getId());
      await ctx.enrollUser(userId, courseIds[i]);
    }
    
    // Lookup the user
    const lookupResp = await ctx.lookupUserById([userId]);
    const user = lookupResp.getResultList()[0];
    
    expect(user.getId()).toBe(userId);
    expect(user.getEnrollmentsList()).toHaveLength(3);
    
    const enrolledCourseIds = user.getEnrollmentsList().map(e => e.getCourse()!.getId());
    courseIds.forEach(id => expect(enrolledCourseIds).toContain(id));
    
    // Verify enrollment details
    const firstEnrollment = user.getEnrollmentsList()[0];
    expect(firstEnrollment.getUser()!.getId()).toBe(userId);
    expect(firstEnrollment.getProgress()).toBe(0);
  });
});


// ============================================================================
// Helper Functions
// ============================================================================

// Test utilities
class TestContext {
  subprocess: Subprocess;
  client: CoursesServiceClient;

  constructor(subprocess: Subprocess, client: CoursesServiceClient) {
    this.subprocess = subprocess;
    this.client = client;
  }

  async cleanup() {
    this.client.close();
    this.subprocess.kill();
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  async queryCourses(): Promise<QueryCoursesResponse> {
    return new Promise((resolve, reject) => {
      const req = new QueryCoursesRequest();
      this.client.queryCourses(req, (err, resp) => {
        if (err) reject(err);
        else if (!resp) reject(new Error("empty response"));
        else resolve(resp);
      });
    });
  }

  async queryCourse(id: string): Promise<QueryCourseResponse> {
    return new Promise((resolve, reject) => {
      const req = new QueryCourseRequest();
      req.setId(id);
      this.client.queryCourse(req, (err, resp) => {
        if (err) reject(err);
        else if (!resp) reject(new Error("empty response"));
        else resolve(resp);
      });
    });
  }

  async createCourse(
    title: string,
    description: string,
    instructorId: string,
    durationHours: number,
    published: boolean
  ): Promise<MutationCreateCourseResponse> {
    return new Promise((resolve, reject) => {
      const input = new CourseInput();
      input.setTitle(title);
      input.setDescription(description);
      input.setInstructorId(instructorId);
      input.setDurationHours(durationHours);
      input.setPublished(published);

      const req = new MutationCreateCourseRequest();
      req.setInput(input);

      this.client.mutationCreateCourse(req, (err, resp) => {
        if (err) reject(err);
        else if (!resp) reject(new Error("empty response"));
        else resolve(resp);
      });
    });
  }

  async enrollUser(userId: string, courseId: string): Promise<MutationEnrollUserResponse> {
    return new Promise((resolve, reject) => {
      const req = new MutationEnrollUserRequest();
      req.setUserId(userId);
      req.setCourseId(courseId);

      this.client.mutationEnrollUser(req, (err, resp) => {
        if (err) reject(err);
        else if (!resp) reject(new Error("empty response"));
        else resolve(resp);
      });
    });
  }

  async lookupUserById(userIds: string[]): Promise<LookupUserByIdResponse> {
    return new Promise((resolve, reject) => {
      const req = new LookupUserByIdRequest();
      const keys = userIds.map(id => {
        const key = new LookupUserByIdRequestKey();
        key.setId(id);
        return key;
      });
      req.setKeysList(keys);

      this.client.lookupUserById(req, (err, resp) => {
        if (err) reject(err);
        else if (!resp) reject(new Error("empty response"));
        else resolve(resp);
      });
    });
  }
}

async function createTestContext(): Promise<TestContext> {
  const proc = Bun.spawn(["bun", "run", "src/plugin.ts"], {
    stdout: "pipe",
    stderr: "inherit",
  });

  if (!proc.stdout) {
    throw new Error("plugin stdout not available");
  }

  const reader = proc.stdout.getReader();
  const decoder = new TextDecoder();
  const { value } = await reader.read();
  reader.releaseLock();

  const text = decoder.decode(value ?? new Uint8Array());
  const firstLine = text.split("\n")[0]?.trim() ?? "";
  const parts = firstLine.split("|");
  const address = parts[3];

  const target = 'unix://' + address;
  const client = new CoursesServiceClient(target, grpc.credentials.createInsecure());

  return new TestContext(proc, client);
}

// Wrapper function that auto-handles test context lifecycle
async function testWithContext(
  name: string,
  testFn: (ctx: TestContext) => Promise<void>
) {
  test(name, async () => {
    const ctx = await createTestContext();
    try {
      await testFn(ctx);
    } finally {
      await ctx.cleanup();
    }
  });
}