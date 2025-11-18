// package: service
// file: service.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as service_pb from "./service_pb";

interface ICoursesServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    lookupUserById: ICoursesServiceService_ILookupUserById;
    mutationCreateCourse: ICoursesServiceService_IMutationCreateCourse;
    mutationEnrollUser: ICoursesServiceService_IMutationEnrollUser;
    queryCourse: ICoursesServiceService_IQueryCourse;
    queryCourses: ICoursesServiceService_IQueryCourses;
}

interface ICoursesServiceService_ILookupUserById extends grpc.MethodDefinition<service_pb.LookupUserByIdRequest, service_pb.LookupUserByIdResponse> {
    path: "/service.CoursesService/LookupUserById";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<service_pb.LookupUserByIdRequest>;
    requestDeserialize: grpc.deserialize<service_pb.LookupUserByIdRequest>;
    responseSerialize: grpc.serialize<service_pb.LookupUserByIdResponse>;
    responseDeserialize: grpc.deserialize<service_pb.LookupUserByIdResponse>;
}
interface ICoursesServiceService_IMutationCreateCourse extends grpc.MethodDefinition<service_pb.MutationCreateCourseRequest, service_pb.MutationCreateCourseResponse> {
    path: "/service.CoursesService/MutationCreateCourse";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<service_pb.MutationCreateCourseRequest>;
    requestDeserialize: grpc.deserialize<service_pb.MutationCreateCourseRequest>;
    responseSerialize: grpc.serialize<service_pb.MutationCreateCourseResponse>;
    responseDeserialize: grpc.deserialize<service_pb.MutationCreateCourseResponse>;
}
interface ICoursesServiceService_IMutationEnrollUser extends grpc.MethodDefinition<service_pb.MutationEnrollUserRequest, service_pb.MutationEnrollUserResponse> {
    path: "/service.CoursesService/MutationEnrollUser";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<service_pb.MutationEnrollUserRequest>;
    requestDeserialize: grpc.deserialize<service_pb.MutationEnrollUserRequest>;
    responseSerialize: grpc.serialize<service_pb.MutationEnrollUserResponse>;
    responseDeserialize: grpc.deserialize<service_pb.MutationEnrollUserResponse>;
}
interface ICoursesServiceService_IQueryCourse extends grpc.MethodDefinition<service_pb.QueryCourseRequest, service_pb.QueryCourseResponse> {
    path: "/service.CoursesService/QueryCourse";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<service_pb.QueryCourseRequest>;
    requestDeserialize: grpc.deserialize<service_pb.QueryCourseRequest>;
    responseSerialize: grpc.serialize<service_pb.QueryCourseResponse>;
    responseDeserialize: grpc.deserialize<service_pb.QueryCourseResponse>;
}
interface ICoursesServiceService_IQueryCourses extends grpc.MethodDefinition<service_pb.QueryCoursesRequest, service_pb.QueryCoursesResponse> {
    path: "/service.CoursesService/QueryCourses";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<service_pb.QueryCoursesRequest>;
    requestDeserialize: grpc.deserialize<service_pb.QueryCoursesRequest>;
    responseSerialize: grpc.serialize<service_pb.QueryCoursesResponse>;
    responseDeserialize: grpc.deserialize<service_pb.QueryCoursesResponse>;
}

export const CoursesServiceService: ICoursesServiceService;

export interface ICoursesServiceServer extends grpc.UntypedServiceImplementation {
    lookupUserById: grpc.handleUnaryCall<service_pb.LookupUserByIdRequest, service_pb.LookupUserByIdResponse>;
    mutationCreateCourse: grpc.handleUnaryCall<service_pb.MutationCreateCourseRequest, service_pb.MutationCreateCourseResponse>;
    mutationEnrollUser: grpc.handleUnaryCall<service_pb.MutationEnrollUserRequest, service_pb.MutationEnrollUserResponse>;
    queryCourse: grpc.handleUnaryCall<service_pb.QueryCourseRequest, service_pb.QueryCourseResponse>;
    queryCourses: grpc.handleUnaryCall<service_pb.QueryCoursesRequest, service_pb.QueryCoursesResponse>;
}

export interface ICoursesServiceClient {
    lookupUserById(request: service_pb.LookupUserByIdRequest, callback: (error: grpc.ServiceError | null, response: service_pb.LookupUserByIdResponse) => void): grpc.ClientUnaryCall;
    lookupUserById(request: service_pb.LookupUserByIdRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.LookupUserByIdResponse) => void): grpc.ClientUnaryCall;
    lookupUserById(request: service_pb.LookupUserByIdRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.LookupUserByIdResponse) => void): grpc.ClientUnaryCall;
    mutationCreateCourse(request: service_pb.MutationCreateCourseRequest, callback: (error: grpc.ServiceError | null, response: service_pb.MutationCreateCourseResponse) => void): grpc.ClientUnaryCall;
    mutationCreateCourse(request: service_pb.MutationCreateCourseRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.MutationCreateCourseResponse) => void): grpc.ClientUnaryCall;
    mutationCreateCourse(request: service_pb.MutationCreateCourseRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.MutationCreateCourseResponse) => void): grpc.ClientUnaryCall;
    mutationEnrollUser(request: service_pb.MutationEnrollUserRequest, callback: (error: grpc.ServiceError | null, response: service_pb.MutationEnrollUserResponse) => void): grpc.ClientUnaryCall;
    mutationEnrollUser(request: service_pb.MutationEnrollUserRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.MutationEnrollUserResponse) => void): grpc.ClientUnaryCall;
    mutationEnrollUser(request: service_pb.MutationEnrollUserRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.MutationEnrollUserResponse) => void): grpc.ClientUnaryCall;
    queryCourse(request: service_pb.QueryCourseRequest, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCourseResponse) => void): grpc.ClientUnaryCall;
    queryCourse(request: service_pb.QueryCourseRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCourseResponse) => void): grpc.ClientUnaryCall;
    queryCourse(request: service_pb.QueryCourseRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCourseResponse) => void): grpc.ClientUnaryCall;
    queryCourses(request: service_pb.QueryCoursesRequest, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCoursesResponse) => void): grpc.ClientUnaryCall;
    queryCourses(request: service_pb.QueryCoursesRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCoursesResponse) => void): grpc.ClientUnaryCall;
    queryCourses(request: service_pb.QueryCoursesRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCoursesResponse) => void): grpc.ClientUnaryCall;
}

export class CoursesServiceClient extends grpc.Client implements ICoursesServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public lookupUserById(request: service_pb.LookupUserByIdRequest, callback: (error: grpc.ServiceError | null, response: service_pb.LookupUserByIdResponse) => void): grpc.ClientUnaryCall;
    public lookupUserById(request: service_pb.LookupUserByIdRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.LookupUserByIdResponse) => void): grpc.ClientUnaryCall;
    public lookupUserById(request: service_pb.LookupUserByIdRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.LookupUserByIdResponse) => void): grpc.ClientUnaryCall;
    public mutationCreateCourse(request: service_pb.MutationCreateCourseRequest, callback: (error: grpc.ServiceError | null, response: service_pb.MutationCreateCourseResponse) => void): grpc.ClientUnaryCall;
    public mutationCreateCourse(request: service_pb.MutationCreateCourseRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.MutationCreateCourseResponse) => void): grpc.ClientUnaryCall;
    public mutationCreateCourse(request: service_pb.MutationCreateCourseRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.MutationCreateCourseResponse) => void): grpc.ClientUnaryCall;
    public mutationEnrollUser(request: service_pb.MutationEnrollUserRequest, callback: (error: grpc.ServiceError | null, response: service_pb.MutationEnrollUserResponse) => void): grpc.ClientUnaryCall;
    public mutationEnrollUser(request: service_pb.MutationEnrollUserRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.MutationEnrollUserResponse) => void): grpc.ClientUnaryCall;
    public mutationEnrollUser(request: service_pb.MutationEnrollUserRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.MutationEnrollUserResponse) => void): grpc.ClientUnaryCall;
    public queryCourse(request: service_pb.QueryCourseRequest, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCourseResponse) => void): grpc.ClientUnaryCall;
    public queryCourse(request: service_pb.QueryCourseRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCourseResponse) => void): grpc.ClientUnaryCall;
    public queryCourse(request: service_pb.QueryCourseRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCourseResponse) => void): grpc.ClientUnaryCall;
    public queryCourses(request: service_pb.QueryCoursesRequest, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCoursesResponse) => void): grpc.ClientUnaryCall;
    public queryCourses(request: service_pb.QueryCoursesRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCoursesResponse) => void): grpc.ClientUnaryCall;
    public queryCourses(request: service_pb.QueryCoursesRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: service_pb.QueryCoursesResponse) => void): grpc.ClientUnaryCall;
}
