import { NEWLINE } from "./const";
import { UtilFunc } from "./types";

export const splitLines = (str: string) => str.split(NEWLINE);

export const makeUtilFunc = (name: string): UtilFunc => ({ name });

export const strBetween = (str: string, left: string, right: string): string =>
    str.split(left)[1].split(right)[0];

export const removeNull = <T = unknown>(array: Array<T | null>): T[] =>
    array.filter((item) => item != null) as T[];
