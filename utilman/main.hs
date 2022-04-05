import System.Environment
import Data.List

main = do
    args <- getArgs
    putStrLn "Args are:"
    mapM putStrLn args

