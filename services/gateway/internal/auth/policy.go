package auth

import (
	"fmt"
	"strings"

	authv1 "go-market/gen/go/auth/v1"

	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Policy struct{}

func (Policy) IsPublic(fullMethod string) (bool, error) {
	// Стандартный Health API не содержит нашу authorization-аннотацию.
	// Метод Check используется только для инфраструктурных проверок готовности.
	if fullMethod == healthv1.Health_Check_FullMethodName {
		return true, nil
	}

	descriptorName, err := methodDescriptorName(fullMethod)
	if err != nil {
		return false, err
	}

	descriptor, err := protoregistry.GlobalFiles.FindDescriptorByName(
		descriptorName,
	)
	if err != nil {
		return false, fmt.Errorf(
			"find descriptor for method %q: %w",
			fullMethod,
			err,
		)
	}

	method, ok := descriptor.(protoreflect.MethodDescriptor)
	if !ok {
		return false, fmt.Errorf(
			"descriptor %q is not a method",
			descriptorName,
		)
	}

	options, ok := method.Options().(*descriptorpb.MethodOptions)
	if !ok {
		return false, fmt.Errorf(
			"read options for method %q",
			fullMethod,
		)
	}

	if !proto.HasExtension(options, authv1.E_Authorization) {
		return false, nil
	}

	value := proto.GetExtension(
		options,
		authv1.E_Authorization,
	)

	rule, ok := value.(*authv1.AuthorizationRule)
	if !ok {
		return false, fmt.Errorf(
			"invalid authorization rule for method %q",
			fullMethod,
		)
	}

	return rule.GetPublic(), nil
}

func methodDescriptorName(
	fullMethod string,
) (protoreflect.FullName, error) {
	value, ok := strings.CutPrefix(fullMethod, "/")
	if !ok {
		return "", fmt.Errorf(
			"invalid gRPC method name %q",
			fullMethod,
		)
	}

	service, method, ok := strings.Cut(value, "/")
	if !ok ||
		service == "" ||
		method == "" ||
		strings.Contains(method, "/") {
		return "", fmt.Errorf(
			"invalid gRPC method name %q",
			fullMethod,
		)
	}

	name := protoreflect.FullName(service + "." + method)
	if !name.IsValid() {
		return "", fmt.Errorf(
			"invalid protobuf method name %q",
			name,
		)
	}

	return name, nil
}
