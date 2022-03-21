import {
    makeFunctionsFromStr,
    makeFunctionFromLine,
    makeFunctionFromLineAnon,
    makeFunctionFromLineKeyword,
} from "../src/main";
import { makeUtilFunc } from "../src/util";

describe("makeFunctionsFromString", () => {
    it("empty string should have no function objects", () => {
        expect(makeFunctionsFromStr("")).toStrictEqual([]);
    });

    it("should return single add function", () => {
        const fileContent = "export const add = (l, r) => l + r;\n";
        const expectedFunctions = [makeUtilFunc("add")];
        expect(makeFunctionsFromStr(fileContent)).toStrictEqual(
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

    it("keyword function should return func obj", () => {
        expect(makeFunctionFromLine("export function fn() {}")).toStrictEqual(
            makeUtilFunc("fn"),
        );
        expect(
            makeFunctionFromLine("export async function fn() {}"),
        ).toStrictEqual(makeUtilFunc("fn"));
    });

    it("constants should return null", () => {
        expect(makeFunctionFromLine('export const FIZZ = "buzz"')).toBeNull();
    });
});

describe("makeFunctionFromLineAnon", () => {
    it("single line anon function should return func obj", () => {
        expect(
            makeFunctionFromLineAnon("export const fn = () => {}"),
        ).toStrictEqual(makeUtilFunc("fn"));
    });
});

describe("makeFunctionFromLineKeyword", () => {
    it("single line function should return func obj", () => {
        expect(
            makeFunctionFromLineKeyword("export function fn() {}"),
        ).toStrictEqual(makeUtilFunc("fn"));
    });
});
