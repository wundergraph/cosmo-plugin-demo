// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var service_pb = require('./service_pb.js');

function serialize_service_LookupUserByIdRequest(arg) {
  if (!(arg instanceof service_pb.LookupUserByIdRequest)) {
    throw new Error('Expected argument of type service.LookupUserByIdRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_LookupUserByIdRequest(buffer_arg) {
  return service_pb.LookupUserByIdRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_LookupUserByIdResponse(arg) {
  if (!(arg instanceof service_pb.LookupUserByIdResponse)) {
    throw new Error('Expected argument of type service.LookupUserByIdResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_LookupUserByIdResponse(buffer_arg) {
  return service_pb.LookupUserByIdResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_MutationCreateCourseRequest(arg) {
  if (!(arg instanceof service_pb.MutationCreateCourseRequest)) {
    throw new Error('Expected argument of type service.MutationCreateCourseRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_MutationCreateCourseRequest(buffer_arg) {
  return service_pb.MutationCreateCourseRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_MutationCreateCourseResponse(arg) {
  if (!(arg instanceof service_pb.MutationCreateCourseResponse)) {
    throw new Error('Expected argument of type service.MutationCreateCourseResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_MutationCreateCourseResponse(buffer_arg) {
  return service_pb.MutationCreateCourseResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_MutationEnrollUserRequest(arg) {
  if (!(arg instanceof service_pb.MutationEnrollUserRequest)) {
    throw new Error('Expected argument of type service.MutationEnrollUserRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_MutationEnrollUserRequest(buffer_arg) {
  return service_pb.MutationEnrollUserRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_MutationEnrollUserResponse(arg) {
  if (!(arg instanceof service_pb.MutationEnrollUserResponse)) {
    throw new Error('Expected argument of type service.MutationEnrollUserResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_MutationEnrollUserResponse(buffer_arg) {
  return service_pb.MutationEnrollUserResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_QueryCourseRequest(arg) {
  if (!(arg instanceof service_pb.QueryCourseRequest)) {
    throw new Error('Expected argument of type service.QueryCourseRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_QueryCourseRequest(buffer_arg) {
  return service_pb.QueryCourseRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_QueryCourseResponse(arg) {
  if (!(arg instanceof service_pb.QueryCourseResponse)) {
    throw new Error('Expected argument of type service.QueryCourseResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_QueryCourseResponse(buffer_arg) {
  return service_pb.QueryCourseResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_QueryCoursesRequest(arg) {
  if (!(arg instanceof service_pb.QueryCoursesRequest)) {
    throw new Error('Expected argument of type service.QueryCoursesRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_QueryCoursesRequest(buffer_arg) {
  return service_pb.QueryCoursesRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_service_QueryCoursesResponse(arg) {
  if (!(arg instanceof service_pb.QueryCoursesResponse)) {
    throw new Error('Expected argument of type service.QueryCoursesResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_service_QueryCoursesResponse(buffer_arg) {
  return service_pb.QueryCoursesResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// Service definition for CoursesService
var CoursesServiceService = exports.CoursesServiceService = {
  // Lookup User entity by id: User type from users plugin (extended here)
lookupUserById: {
    path: '/service.CoursesService/LookupUserById',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.LookupUserByIdRequest,
    responseType: service_pb.LookupUserByIdResponse,
    requestSerialize: serialize_service_LookupUserByIdRequest,
    requestDeserialize: deserialize_service_LookupUserByIdRequest,
    responseSerialize: serialize_service_LookupUserByIdResponse,
    responseDeserialize: deserialize_service_LookupUserByIdResponse,
  },
  // Creates a new course
mutationCreateCourse: {
    path: '/service.CoursesService/MutationCreateCourse',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.MutationCreateCourseRequest,
    responseType: service_pb.MutationCreateCourseResponse,
    requestSerialize: serialize_service_MutationCreateCourseRequest,
    requestDeserialize: deserialize_service_MutationCreateCourseRequest,
    responseSerialize: serialize_service_MutationCreateCourseResponse,
    responseDeserialize: deserialize_service_MutationCreateCourseResponse,
  },
  // Enrolls a user in a course
mutationEnrollUser: {
    path: '/service.CoursesService/MutationEnrollUser',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.MutationEnrollUserRequest,
    responseType: service_pb.MutationEnrollUserResponse,
    requestSerialize: serialize_service_MutationEnrollUserRequest,
    requestDeserialize: deserialize_service_MutationEnrollUserRequest,
    responseSerialize: serialize_service_MutationEnrollUserResponse,
    responseDeserialize: deserialize_service_MutationEnrollUserResponse,
  },
  // Returns a single course by ID
queryCourse: {
    path: '/service.CoursesService/QueryCourse',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.QueryCourseRequest,
    responseType: service_pb.QueryCourseResponse,
    requestSerialize: serialize_service_QueryCourseRequest,
    requestDeserialize: deserialize_service_QueryCourseRequest,
    responseSerialize: serialize_service_QueryCourseResponse,
    responseDeserialize: deserialize_service_QueryCourseResponse,
  },
  // Returns a list of all courses
queryCourses: {
    path: '/service.CoursesService/QueryCourses',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.QueryCoursesRequest,
    responseType: service_pb.QueryCoursesResponse,
    requestSerialize: serialize_service_QueryCoursesRequest,
    requestDeserialize: deserialize_service_QueryCoursesRequest,
    responseSerialize: serialize_service_QueryCoursesResponse,
    responseDeserialize: deserialize_service_QueryCoursesResponse,
  },
};

exports.CoursesServiceClient = grpc.makeGenericClientConstructor(CoursesServiceService, 'CoursesService');
