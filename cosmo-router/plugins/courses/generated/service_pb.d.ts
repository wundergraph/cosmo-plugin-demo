// package: service
// file: service.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class LookupUserByIdRequestKey extends jspb.Message { 
    getId(): string;
    setId(value: string): LookupUserByIdRequestKey;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LookupUserByIdRequestKey.AsObject;
    static toObject(includeInstance: boolean, msg: LookupUserByIdRequestKey): LookupUserByIdRequestKey.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LookupUserByIdRequestKey, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LookupUserByIdRequestKey;
    static deserializeBinaryFromReader(message: LookupUserByIdRequestKey, reader: jspb.BinaryReader): LookupUserByIdRequestKey;
}

export namespace LookupUserByIdRequestKey {
    export type AsObject = {
        id: string,
    }
}

export class LookupUserByIdRequest extends jspb.Message { 
    clearKeysList(): void;
    getKeysList(): Array<LookupUserByIdRequestKey>;
    setKeysList(value: Array<LookupUserByIdRequestKey>): LookupUserByIdRequest;
    addKeys(value?: LookupUserByIdRequestKey, index?: number): LookupUserByIdRequestKey;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LookupUserByIdRequest.AsObject;
    static toObject(includeInstance: boolean, msg: LookupUserByIdRequest): LookupUserByIdRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LookupUserByIdRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LookupUserByIdRequest;
    static deserializeBinaryFromReader(message: LookupUserByIdRequest, reader: jspb.BinaryReader): LookupUserByIdRequest;
}

export namespace LookupUserByIdRequest {
    export type AsObject = {
        keysList: Array<LookupUserByIdRequestKey.AsObject>,
    }
}

export class LookupUserByIdResponse extends jspb.Message { 
    clearResultList(): void;
    getResultList(): Array<User>;
    setResultList(value: Array<User>): LookupUserByIdResponse;
    addResult(value?: User, index?: number): User;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LookupUserByIdResponse.AsObject;
    static toObject(includeInstance: boolean, msg: LookupUserByIdResponse): LookupUserByIdResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LookupUserByIdResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LookupUserByIdResponse;
    static deserializeBinaryFromReader(message: LookupUserByIdResponse, reader: jspb.BinaryReader): LookupUserByIdResponse;
}

export namespace LookupUserByIdResponse {
    export type AsObject = {
        resultList: Array<User.AsObject>,
    }
}

export class QueryCoursesRequest extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryCoursesRequest.AsObject;
    static toObject(includeInstance: boolean, msg: QueryCoursesRequest): QueryCoursesRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryCoursesRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryCoursesRequest;
    static deserializeBinaryFromReader(message: QueryCoursesRequest, reader: jspb.BinaryReader): QueryCoursesRequest;
}

export namespace QueryCoursesRequest {
    export type AsObject = {
    }
}

export class QueryCoursesResponse extends jspb.Message { 
    clearCoursesList(): void;
    getCoursesList(): Array<Course>;
    setCoursesList(value: Array<Course>): QueryCoursesResponse;
    addCourses(value?: Course, index?: number): Course;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryCoursesResponse.AsObject;
    static toObject(includeInstance: boolean, msg: QueryCoursesResponse): QueryCoursesResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryCoursesResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryCoursesResponse;
    static deserializeBinaryFromReader(message: QueryCoursesResponse, reader: jspb.BinaryReader): QueryCoursesResponse;
}

export namespace QueryCoursesResponse {
    export type AsObject = {
        coursesList: Array<Course.AsObject>,
    }
}

export class QueryCourseRequest extends jspb.Message { 
    getId(): string;
    setId(value: string): QueryCourseRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryCourseRequest.AsObject;
    static toObject(includeInstance: boolean, msg: QueryCourseRequest): QueryCourseRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryCourseRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryCourseRequest;
    static deserializeBinaryFromReader(message: QueryCourseRequest, reader: jspb.BinaryReader): QueryCourseRequest;
}

export namespace QueryCourseRequest {
    export type AsObject = {
        id: string,
    }
}

export class QueryCourseResponse extends jspb.Message { 

    hasCourse(): boolean;
    clearCourse(): void;
    getCourse(): Course | undefined;
    setCourse(value?: Course): QueryCourseResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryCourseResponse.AsObject;
    static toObject(includeInstance: boolean, msg: QueryCourseResponse): QueryCourseResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryCourseResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryCourseResponse;
    static deserializeBinaryFromReader(message: QueryCourseResponse, reader: jspb.BinaryReader): QueryCourseResponse;
}

export namespace QueryCourseResponse {
    export type AsObject = {
        course?: Course.AsObject,
    }
}

export class MutationCreateCourseRequest extends jspb.Message { 

    hasInput(): boolean;
    clearInput(): void;
    getInput(): CourseInput | undefined;
    setInput(value?: CourseInput): MutationCreateCourseRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MutationCreateCourseRequest.AsObject;
    static toObject(includeInstance: boolean, msg: MutationCreateCourseRequest): MutationCreateCourseRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MutationCreateCourseRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MutationCreateCourseRequest;
    static deserializeBinaryFromReader(message: MutationCreateCourseRequest, reader: jspb.BinaryReader): MutationCreateCourseRequest;
}

export namespace MutationCreateCourseRequest {
    export type AsObject = {
        input?: CourseInput.AsObject,
    }
}

export class MutationCreateCourseResponse extends jspb.Message { 

    hasCreateCourse(): boolean;
    clearCreateCourse(): void;
    getCreateCourse(): Course | undefined;
    setCreateCourse(value?: Course): MutationCreateCourseResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MutationCreateCourseResponse.AsObject;
    static toObject(includeInstance: boolean, msg: MutationCreateCourseResponse): MutationCreateCourseResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MutationCreateCourseResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MutationCreateCourseResponse;
    static deserializeBinaryFromReader(message: MutationCreateCourseResponse, reader: jspb.BinaryReader): MutationCreateCourseResponse;
}

export namespace MutationCreateCourseResponse {
    export type AsObject = {
        createCourse?: Course.AsObject,
    }
}

export class MutationEnrollUserRequest extends jspb.Message { 
    getUserId(): string;
    setUserId(value: string): MutationEnrollUserRequest;
    getCourseId(): string;
    setCourseId(value: string): MutationEnrollUserRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MutationEnrollUserRequest.AsObject;
    static toObject(includeInstance: boolean, msg: MutationEnrollUserRequest): MutationEnrollUserRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MutationEnrollUserRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MutationEnrollUserRequest;
    static deserializeBinaryFromReader(message: MutationEnrollUserRequest, reader: jspb.BinaryReader): MutationEnrollUserRequest;
}

export namespace MutationEnrollUserRequest {
    export type AsObject = {
        userId: string,
        courseId: string,
    }
}

export class MutationEnrollUserResponse extends jspb.Message { 

    hasEnrollUser(): boolean;
    clearEnrollUser(): void;
    getEnrollUser(): Enrollment | undefined;
    setEnrollUser(value?: Enrollment): MutationEnrollUserResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MutationEnrollUserResponse.AsObject;
    static toObject(includeInstance: boolean, msg: MutationEnrollUserResponse): MutationEnrollUserResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MutationEnrollUserResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MutationEnrollUserResponse;
    static deserializeBinaryFromReader(message: MutationEnrollUserResponse, reader: jspb.BinaryReader): MutationEnrollUserResponse;
}

export namespace MutationEnrollUserResponse {
    export type AsObject = {
        enrollUser?: Enrollment.AsObject,
    }
}

export class User extends jspb.Message { 
    getId(): string;
    setId(value: string): User;
    clearInstructorCoursesList(): void;
    getInstructorCoursesList(): Array<Course>;
    setInstructorCoursesList(value: Array<Course>): User;
    addInstructorCourses(value?: Course, index?: number): Course;
    clearEnrollmentsList(): void;
    getEnrollmentsList(): Array<Enrollment>;
    setEnrollmentsList(value: Array<Enrollment>): User;
    addEnrollments(value?: Enrollment, index?: number): Enrollment;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): User.AsObject;
    static toObject(includeInstance: boolean, msg: User): User.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): User;
    static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
    export type AsObject = {
        id: string,
        instructorCoursesList: Array<Course.AsObject>,
        enrollmentsList: Array<Enrollment.AsObject>,
    }
}

export class Course extends jspb.Message { 
    getId(): string;
    setId(value: string): Course;
    getTitle(): string;
    setTitle(value: string): Course;
    getDescription(): string;
    setDescription(value: string): Course;

    hasInstructor(): boolean;
    clearInstructor(): void;
    getInstructor(): User | undefined;
    setInstructor(value?: User): Course;
    getDurationHours(): number;
    setDurationHours(value: number): Course;
    getPublished(): boolean;
    setPublished(value: boolean): Course;
    clearEnrollmentsList(): void;
    getEnrollmentsList(): Array<Enrollment>;
    setEnrollmentsList(value: Array<Enrollment>): Course;
    addEnrollments(value?: Enrollment, index?: number): Enrollment;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Course.AsObject;
    static toObject(includeInstance: boolean, msg: Course): Course.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Course, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Course;
    static deserializeBinaryFromReader(message: Course, reader: jspb.BinaryReader): Course;
}

export namespace Course {
    export type AsObject = {
        id: string,
        title: string,
        description: string,
        instructor?: User.AsObject,
        durationHours: number,
        published: boolean,
        enrollmentsList: Array<Enrollment.AsObject>,
    }
}

export class CourseInput extends jspb.Message { 
    getTitle(): string;
    setTitle(value: string): CourseInput;
    getDescription(): string;
    setDescription(value: string): CourseInput;
    getInstructorId(): string;
    setInstructorId(value: string): CourseInput;
    getDurationHours(): number;
    setDurationHours(value: number): CourseInput;
    getPublished(): boolean;
    setPublished(value: boolean): CourseInput;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CourseInput.AsObject;
    static toObject(includeInstance: boolean, msg: CourseInput): CourseInput.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CourseInput, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CourseInput;
    static deserializeBinaryFromReader(message: CourseInput, reader: jspb.BinaryReader): CourseInput;
}

export namespace CourseInput {
    export type AsObject = {
        title: string,
        description: string,
        instructorId: string,
        durationHours: number,
        published: boolean,
    }
}

export class Enrollment extends jspb.Message { 
    getId(): string;
    setId(value: string): Enrollment;

    hasUser(): boolean;
    clearUser(): void;
    getUser(): User | undefined;
    setUser(value?: User): Enrollment;

    hasCourse(): boolean;
    clearCourse(): void;
    getCourse(): Course | undefined;
    setCourse(value?: Course): Enrollment;
    getProgress(): number;
    setProgress(value: number): Enrollment;
    getEnrolledAt(): string;
    setEnrolledAt(value: string): Enrollment;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Enrollment.AsObject;
    static toObject(includeInstance: boolean, msg: Enrollment): Enrollment.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Enrollment, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Enrollment;
    static deserializeBinaryFromReader(message: Enrollment, reader: jspb.BinaryReader): Enrollment;
}

export namespace Enrollment {
    export type AsObject = {
        id: string,
        user?: User.AsObject,
        course?: Course.AsObject,
        progress: number,
        enrolledAt: string,
    }
}
