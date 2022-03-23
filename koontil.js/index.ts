// This code is imported from k2on/koontiljs (Oct 12, 2020)

const colors = require('colors');


export const makeTable = (table: { [key: string]: any; }) => {
    let longestKeyLength = Math.max(...Object.keys(table).map(s => s.length as number));
    let tableString = '';
    for (let key in table) {
        let value = table[key];
        let padding = "  ";
        for (let i=0;i<longestKeyLength-key.length;i++) {
            padding += " ";
        }
        tableString += `${key}:${padding}${value}\n`;
    } 
    return tableString;
}

export const formatString = (string: string, obj: { [key: string]: any; }): string => {
    for (let key in obj) {
        let value = obj[key];
        string = string.replace(new RegExp(`{${key}}`, "g"), value);
    }
    return string;
}

export const titleCase = (str: string) => str.replace(
    /\w\S*/g,
    function(txt) {
        return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    }
);

export const prettyPrint = (v: any) => console.log(prettyFormat(v));

const makeTabs = (tabs: number) => {
    let s = '';
    for (let i=0;i<tabs;i++) {
        s += '  ';
    }
    return s;
}

const prettyFormat = (v: any, tabIndex = 0, lastElementList=false): string => {
    // first few are simple
    if (typeof v == 'string') return colors.green(v);
    if (typeof v == 'number') return colors.cyan(v);
    if (typeof v == 'boolean') return colors.blue(v);
    if (typeof v == 'undefined') return colors.grey('undefined');
    if (v == null) return colors.grey('null');
    if (v instanceof Set) Array.from(v);
    
    if (Array.isArray(v)) return v.map(i => `\n${makeTabs(tabIndex)}- ${prettyFormat(i, tabIndex+1, true)}`).join('')

    if (typeof v == 'object') {
        let s = '';
        let i = 0;
        for (let key in v) {
            s += `${i == 0 && lastElementList ? '' : `\n${makeTabs(tabIndex)}`}${key}: ${prettyFormat(v[key], tabIndex+1)}`;
            i++;
        }
        return s;
    }
    throw new Error('unsupported pretty print');
}
