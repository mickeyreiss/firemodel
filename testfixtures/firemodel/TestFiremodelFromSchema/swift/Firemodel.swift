// DO NOT EDIT - Code generated by firemodel (dev).

import Foundation
import Pring

// TODO: Add documentation to TestEnum in firemodel schema.
@objc enum TestEnum: Int {
    // TODO: Add documentation to Left in firemodel schema.
    case left
    // TODO: Add documentation to Right in firemodel schema.
    case right
    // TODO: Add documentation to Up in firemodel schema.
    case up
    // TODO: Add documentation to Down in firemodel schema.
    case down
}

extension TestEnum: CustomDebugStringConvertible {
    init?(firestoreValue value: Any?) {
        guard let value = value as? String else {
            return nil
        }
        switch value {
        case "LEFT":
            self = .left
        case "RIGHT":
            self = .right
        case "UP":
            self = .up
        case "DOWN":
            self = .down
        default:
            return nil
        }
    }

    var firestoreValue: String? {
        switch self {
        case .left:
            return "LEFT"
        case .right:
            return "RIGHT"
        case .up:
            return "UP"
        case .down:
            return "DOWN"
        }
    }

    var debugDescription: String { return firestoreValue ?? "<INVALID>" }
}

// TODO: Add documentation to TestStruct in firemodel schema.
@objcMembers class TestStruct: Pring.Object {
    // TODO: Add documentation to where in firemodel schema.
    var where: String?
    // TODO: Add documentation to how_much in firemodel schema.
    var howMuch: Int = 0
    // TODO: Add documentation to some_enum in firemodel schema.
    var someEnum: TestEnum?

    override func encode(_ key: String, value: Any?) -> Any? {
        switch key {
        case "someEnum":
            return self.someEnum?.firestoreValue
        default:
            break
        }
        return nil
    }

    override func decode(_ key: String, value: Any?) -> Bool {
        switch key {
        case "someEnum":
            self.someEnum = TestEnum(firestoreValue: value)
        default:
            break
        }
        return false
    }
}

// A Test is a test model.
@objcMembers class TestModel: Pring.Object {
override class var path: String { return "test_models" }
    // The name.
    var name: String?
    // The age.
    var age: Int = 0
    // The number pi.
    var pi: Float = 0
    // The birth date.
    var birthdate: Date?
    // True if it is good.
    var isGood: Bool = false
    // TODO: Add documentation to data in firemodel schema.
    var data: Data?
    // TODO: Add documentation to friend in firemodel schema.
    var friend: Pring.Reference<TestModel> = .init()
    // TODO: Add documentation to location in firemodel schema.
    var location: Pring.GeoPoint?
    // TODO: Add documentation to colors in firemodel schema.
    var colors: [String]?
    // TODO: Add documentation to numbers in firemodel schema.
    var numbers: [Int]?
    // TODO: Add documentation to bools in firemodel schema.
    var bools: [Bool]?
    // TODO: Add documentation to doubles in firemodel schema.
    var doubles: [Float]?
    // TODO: Add documentation to directions in firemodel schema.
    var directions: [TestEnum]?
    // TODO: Add documentation to models in firemodel schema.
    var models: [TestStruct]?
    // TODO: Add documentation to models2 in firemodel schema.
    var models2: [TestStruct]?
    // TODO: Add documentation to refs in firemodel schema.
    var refs: [Pring.AnyReference] = .init()
    // TODO: Add documentation to modelRefs in firemodel schema.
    var modelRefs: [Any] = .init()
    // TODO: Add documentation to meta in firemodel schema.
    var meta: [String: Any] = [:]
    // TODO: Add documentation to metaStrs in firemodel schema.
    var metaStrs: [String: String] = [:]
    // TODO: Add documentation to direction in firemodel schema.
    var direction: TestEnum?
    // TODO: Add documentation to testFile in firemodel schema.
    var testFile: Pring.File?
    // TODO: Add documentation to url in firemodel schema.
    var url: URL?
    // TODO: Add documentation to nested in firemodel schema.
    var nested: TestStruct?
    // TODO: Add documentation to nested_collection in firemodel schema.
    var nestedCollection: Pring.NestedCollection<TestModel> = []

    override func encode(_ key: String, value: Any?) -> Any? {
        switch key {
        case "direction":
            return self.direction?.firestoreValue
        case "models":
            return self.models?.map { $0.rawValue }
        case "models2":
            return self.models2?.map { $0.rawValue }
        case "nested":
            return self.nested?.rawValue
        case "directions":
            return self.directions?.map { $0.firestoreValue }
        default:
            break
        }
        return nil
    }

    override func decode(_ key: String, value: Any?) -> Bool {
        switch key {
        case "direction":
            self.direction = TestEnum(firestoreValue: value)
        case "models":
            self.models = (value as? [[String: Any]])?
                .enumerated()
                .map { TestStruct(id: "models.\($0.offset)", value: $0.element) }
        case "models2":
            self.models2 = (value as? [[String: Any]])?
                .enumerated()
                .map { TestStruct(id: "models2.\($0.offset)", value: $0.element) }
        case "nested":
          if let value = value as? [String: Any] {
            self.nested = TestStruct(id: "\(0)", value: value)
            return true
          }
        case "directions":
            self.directions = (value as? [String])?.compactMap { TestEnum(firestoreValue: $0) }
			return true
        default:
            break
        }
        return false
    }
}

// TODO: Add documentation to TestTimestamps in firemodel schema.
@objcMembers class TestTimestamps: Pring.Object {
override class var path: String { return "timestamps" }
}

// TODO: Add documentation to Test in firemodel schema.
@objcMembers class Test: Pring.Object {
    // TODO: Add documentation to direction in firemodel schema.
    var direction: TestEnum?

    override func encode(_ key: String, value: Any?) -> Any? {
        switch key {
        case "direction":
            return self.direction?.firestoreValue
        default:
            break
        }
        return nil
    }

    override func decode(_ key: String, value: Any?) -> Bool {
        switch key {
        case "direction":
            self.direction = TestEnum(firestoreValue: value)
        default:
            break
        }
        return false
    }
}
