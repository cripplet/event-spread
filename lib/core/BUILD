load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary", "go_test")

go_library(
    name = "spread",
    srcs = ["spread.go"],
    importpath = "github.com/cripplet/event-spread/core/spread",
    deps = [
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
        "@com_github_golang_protobuf//ptypes:go_default_library",
    ],    
)
