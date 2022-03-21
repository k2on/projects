import { UtilFunc } from "./types";
import { makeUtilFunc, removeNull, splitLines, strBetween } from "./util";

export const makeFunctionsFromString = (fileContent: string): UtilFunc[] => {
    const lines = splitLines(fileContent);
    return removeNull<UtilFunc>(lines.map(makeFunctionFromLine));
};

export const makeFunctionFromLine = (line: string): UtilFunc | null => {
    const lineTrimmed = line.trim();
    if (lineTrimmed == "") return null;

    const isLineExporting = lineTrimmed.startsWith("export ");
    if (!isLineExporting) return null;

    const REGEX_FN_ANON = /const \S* = (async )?\(/;
    const isLineAnonFn = lineTrimmed.match(REGEX_FN_ANON);
    if (isLineAnonFn) return makeFunctionFromLineAnon(lineTrimmed);
    return null;
};

export const makeFunctionFromLineAnon = (line: string): UtilFunc => {
    const fnName = strBetween(line, "export const ", " = ");
    return makeUtilFunc(fnName);
};
