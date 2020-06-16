load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary",)

go_library(
    name = "spread",
    srcs = ["spread.go"],
    importpath = "github.com/cripplet/event_spreader/core/spread",
    deps = [
        "//lib/proto:event_spread_go_proto",
	"@org_golang_google_grpc//:go_default_library",
	"@org_golang_google_grpc//status:go_default_library",
	"@org_golang_google_grpc//codes:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library",
    ]
)

