
module Util where

import Data.Text (Text, pack, unpack, stripPrefix, breakOn)
import Data.Maybe (fromMaybe)

subStr :: 
    (Either String Text , Either String Text) ->
    Text ->
    Text
subStr (start, end) full =
    fromMaybe
    (pack [])
    (
        retain (stripPrefix <*>) snd start full
            >>= retain (Just .) fst end
    )
    where retain sub part delim t =
            either
                (Just . const t)
                (sub $ part . flip breakOn t)
                delim