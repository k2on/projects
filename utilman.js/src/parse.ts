import { REGEX_FN_ANON, REGEX_FN_KEYWORD } from "./parse/regex";
import { UtilFunc } from "./parse/types";
import { makeUtilFunc, removeNull, splitLines, strBetween } from "./parse/util";

export const makeFunctionsFromStr = (str: string): UtilFunc[] => {
    const lines = splitLines(str);
    return removeNull<UtilFunc>(lines.map(makeFunctionFromLine));
};

export const makeFunctionFromLine = (
    lineUntrimmed: string,
): UtilFunc | null => {
    const line = lineUntrimmed.trim();
    if (line == "") return null;
    if (!isLineExporting(line)) return null;
    if (isLineAnonFn(line)) return makeFunctionFromLineAnon(line);
    if (isLineKeywordFn(line)) return makeFunctionFromLineKeyword(line);
    return null;
};

const isLineExporting = (line: string) => line.startsWith("export ");
const isLineAnonFn = (line: string) => line.match(REGEX_FN_ANON);
const isLineKeywordFn = (line: string) => line.match(REGEX_FN_KEYWORD);

export const makeFunctionFromLineAnon = (line: string): UtilFunc => {
    const fnName = strBetween(line, "export const ", " = ").trim();
    return makeUtilFunc(fnName);
};

export const makeFunctionFromLineKeyword = (line: string): UtilFunc => {
    const fnName = strBetween(line, " function ", "(").trim();
    return makeUtilFunc(fnName);
};
