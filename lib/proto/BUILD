load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library", "go_grpc_library")

package(default_visibility=["//visibility:public"])

proto_library(
    name = "event_spread_proto",
    srcs = ["event_spread.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_grpc_library(
    name = "event_spread_go_proto",
    importpath = "github.com/cripplet/event_spread/lib/proto/event_spread_go_proto",
    proto = ":event_spread_proto",
)
