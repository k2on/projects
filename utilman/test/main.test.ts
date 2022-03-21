import {
    makeFunctionsFromString,
    makeFunctionFromLine,
    makeFunctionFromLineAnon,
} from "../src/main";
import { makeUtilFunc } from "../src/util";

describe("makeFunctionsFromFileContent", () => {
    it("empty string should have no function objects", () => {
        expect(makeFunctionsFromString("")).toStrictEqual([]);
    });

    it("should return single add function", () => {
        const fileContent = "export const add = (l, r) => l + r;\n";
        const expectedFunctions = [makeUtilFunc("add")];
        expect(makeFunctionsFromString(fileContent)).toStrictEqual(
            expectedFunctions,
        );
    });
});

describe("makeFunctionFromLine", () => {
    it("empty string should not return a function object", () => {
        expect(makeFunctionFromLine("")).toStrictEqual(null);
    });

    it("line that don't start with export should return null", () => {
        expect(makeFunctionFromLine("const myFunc = () => {}")).toBeNull();
        expect(makeFunctionFromLine("function myFunc() {}")).toBeNull();
        expect(makeFunctionFromLine("async function myFunc() {}")).toBeNull();
    });

    it("anonymous functions should return func obj", () => {
        expect(
            makeFunctionFromLine("export const fn = () => {}"),
        ).toStrictEqual(makeUtilFunc("fn"));
        expect(
            makeFunctionFromLine("export const myFunc = async () => {}"),
        ).toStrictEqual(makeUtilFunc("myFunc"));
    });
});

describe("makeFunctionFromLineAnon", () => {
    it("single line anon function should return func obj", () => {
        expect(
            makeFunctionFromLineAnon("export const fn = () => {}"),
        ).toStrictEqual(makeUtilFunc("fn"));
    });
});
