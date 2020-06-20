load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary", "go_test")

go_library(
    name = "spread",
    srcs = ["spread.go"],
    importpath = "github.com/cripplet/event-spread/core/spread",
    deps = [
	":handlers",
        "//lib/proto:event_spread_go_proto",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)

go_test(
    name = "spread_test",
    srcs = ["spread_test.go"],
    embed = [":spread"],
    importpath = "github.com/cripplet/event-spread/core/spread_test",
    deps = [
	":handlers",
        "@com_github_golang_protobuf//ptypes:go_default_library",
        # "@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
    ],    
)

go_library(
    name = "handlers",
    srcs = ["handlers.go"],
    importpath = "github.com/cripplet/event-spread/core/handlers",
    deps = [
        "//lib/proto:event_spread_go_proto",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library",
        "@com_github_golang_protobuf//ptypes/timestamp:go_default_library",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)

go_test(
    name = "handlers_test",
    srcs = ["handlers_test.go"],
    importpath = "github.com/cripplet/event-spread/core/handlers_test",
    embed = [":handlers"],
    deps = [
        "//lib/proto:event_spread_go_proto",
        "@com_github_golang_protobuf//ptypes:go_default_library",
        # "@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
    ],
)
