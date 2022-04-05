module Parse where

import Types
import Util (subStr)
import Data.Maybe (catMaybes)
import Data.Text (pack, unpack)

makeFuncsFromStr :: String -> [UtilFunc]
makeFuncsFromStr "" = []
makeFuncsFromStr str =
    let strLines = lines str
        maybeFuncs = map makeFuncFromLine strLines
    in catMaybes maybeFuncs

makeFuncFromLine :: String -> Maybe UtilFunc
makeFuncFromLine "" = Nothing
-- makeFuncFromLine line = 
--     let x = 4
--     in x

makeFuncFromLineKeyword :: String -> UtilFunc
makeFuncFromLineKeyword line =
    let fnName = subStr (Left "export function ", Right (pack "(")) (pack line) 
    in UtilFunc (unpack fnName)