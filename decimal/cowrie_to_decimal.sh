# Convert cowrie db all decimal fields stored using mongodb 3.4 decimal128
# Before mongodb 3.4, mongodb can not store decimal, stored as string.

db=cowrie

set -e

convert() {
  coll=$1
  shift
  for f in $*; do
    ./mongo_to_decimal.py --db $db --coll $coll --field $f
  done;
}

convert contractBatch creditline repaymoney repaycapitalmoney repayinterestmoney repayfinemoney
convert creditAddPolicy details#percentage
convert creditObject creditlineeach
convert chickenOrder price money
convert goods price
convert goodsOrder sum loans lines#quantity lines#price lines#money
convert settlement sum lines#quantity lines#price lines#money
convert payment dayrate businessmoney paymoney ratemoney settlemoney
convert project coveragerate
