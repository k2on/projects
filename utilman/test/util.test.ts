import { splitLines, makeUtilFunc, strBetween, removeNull } from "../src/util";

describe("splitLines", () => {
    it("empty str should return array of single empty string", () => {
        expect(splitLines("")).toStrictEqual([""]);
    });

    it("str with newline should be array of two elements", () => {
        expect(splitLines("max\nkoon")).toStrictEqual(["max", "koon"]);
    });
});

describe("makeUtilFunc", () => {
    it("should return an object of a basic util function", () => {
        expect(makeUtilFunc("myFunc")).toStrictEqual({ name: "myFunc" });
    });
});

describe("strBetween", () => {
    it("should return the string between two other strings", () => {
        expect(strBetween("hamburger", "ham", "ger")).toEqual("bur");
        expect(
            strBetween(
                "export const myFunName = () => {}",
                "export const ",
                " = ",
            ),
        ).toEqual("myFunName");
    });
});

describe("removeNull", () => {
    it("if given an empty array should return an empty array", () => {
        expect(removeNull([])).toStrictEqual([]);
    });

    it("single null should return an empty array", () => {
        expect(removeNull([null])).toStrictEqual([]);
    });

    it("single int should return input", () => {
        expect(removeNull([1])).toStrictEqual([1]);
    });

    it("remove only null elements", () => {
        expect(removeNull([1, 2, null, 3, 4, null, null, 5])).toStrictEqual([
            1, 2, 3, 4, 5,
        ]);
    });
});
