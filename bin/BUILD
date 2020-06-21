load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "main",
    srcs = ["main.go"],
    importpath = "github.com/cripplet/event-spread/bin/main",
    deps = [
        "//lib/core:handlers",
        "//lib/core:spread",
	"//lib/proto:event_spread_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//reflection:go_default_library",
    ],
)
