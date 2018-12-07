import qualified Data.Text    as Text
import qualified Data.Text.IO as Text
import Data.Char
import Data.Bits
import Data.List.Unique


part1 = do
    ls <- fmap Text.lines (Text.readFile "input.txt")
    let dt = fmap Text.unpack ls
    let occ = map occurrences dt
    let p23 = map twothree occ
    let p11 = map pair_reduce_or p23
    let tup = foldl tuple_add (0,0) p11
    let res = fst(tup) * snd(tup)
    print res

twothree e = [ ((fst(x)+1) `mod` 2, fst(x) `mod` 2) | x <- e, fst(x)==2 || fst(x)==3 ]

tuple_or  (a,b) (c,d) = ((.|.) a b, (.|.) c d)
tuple_add (a,b) (c,d) = (a + c, b + d)

pair_reduce_or :: [(Int, Int)] -> (Int, Int)
pair_reduce_or [] = (0,0)
pair_reduce_or (x:[]) = x
pair_reduce_or (x:xs) = tuple_or x (pair_reduce_or xs)

---

part2 = do
    ls <- fmap Text.lines (Text.readFile "input.txt")
    let a = fmap Text.unpack ls
    let r = [ zip x y | x <- a, y <- a ]
    let c = map diff r
    let b = [ no_space x | x <- c, (length . filter (==' ')) x == 1 ]
    print b

diff :: [(Char, Char)] -> [Char]
diff [] = []
diff (x:xs) = (if fst(x)==snd(x) then [fst(x)] else [' ']) ++ diff xs 

no_space r = [ x | x <- r, x /= ' ' ]
