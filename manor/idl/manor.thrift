namespace go manor.rpc

// TypeRef 结构
struct TypeRef {
    1: string Kind,
    2: string Name,
    3: optional TypeRef OfType
}

// Argument 结构
struct ArgumentNode {
    1: string Name,
    2: TypeRef Type,
    3: optional string DefaultValue
}

// Field 结构
struct FieldNode {
    1: string Name,
    2: optional string Description,
    3: map<string, ArgumentNode> Args,
    4: TypeRef Type,
    5: bool IsDeprecated,
    6: optional string DeprecationReason
}

// EnumValue 结构
struct EnumValueNode {
    1: string Name,
    2: optional string Description,
    3: bool IsDeprecated,
    4: optional string DeprecationReason
}

// 各种节点类型
struct ScalarNode {
    1: string Name,
    2: optional string Description,
    3: bool IsMain
}

struct ObjectNode {
    1: string Name,
    2: optional string Description,
    3: map<string, FieldNode> Fields,
    4: list<string> InterfaceNames,
    5: bool IsModel,
    6: list<string> Scopes,
    7: string Table,
    8: bool IsMain
}

struct InterfaceNode {
    1: string Name,
    2: optional string Description,
    3: map<string, FieldNode> Fields,
    4: bool IsMain
}

struct UnionNode {
    1: string Name,
    2: optional string Description,
    3: map<string, string> TypeNames,
    4: bool IsMain
}

struct EnumNode {
    1: string Name,
    2: optional string Description,
    3: map<string, EnumValueNode> EnumValues,
    4: bool IsMain
}

struct InputObjectNode {
    1: string Name,
    2: optional string Description,
    3: map<string, FieldNode> Fields,
    4: bool IsMain
}

// NodeStore 结构
struct NodeStore {
    1: map<string, ScalarNode> Scalars,
    2: map<string, InterfaceNode> Interfaces,
    3: map<string, ObjectNode> Objects,
    4: map<string, UnionNode> Unions,
    5: map<string, EnumNode> Enums,
    6: map<string, InputObjectNode> Inputs
}

struct RegisterRequest {
    1: string ServiceName,
    2: string ServiceAddr,
    3: NodeStore Store,
}

struct RegisterResponse {
    1: bool Success,
    2: string Message,
}

struct PingRequest {
    1: string ServiceName,
    2: string ServiceAddr,
}

struct PingResponse {
    1: string Message,
}

service Manor {
    RegisterResponse Register(1: RegisterRequest req)
    PingResponse Ping(1: PingRequest req)
}