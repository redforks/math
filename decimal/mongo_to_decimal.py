#!/usr/bin/python3

# Convert a mongodb field from string to decimal128

import argparse
from pymongo import MongoClient
from decimal import Decimal
from bson.decimal128 import Decimal128


def convert_value(v):
    if v is None:
        return ('n', False)

    if type(v) is not str:
        return ('s', False)

    return Decimal128(Decimal(v)), True


def convert_details(v, f):
    if v is None:
        return ('n', False)

    for rec in v:
        val, need = convert_value(rec[f])
        if need:
            rec[f] = val
    return (v, True)


parser = argparse.ArgumentParser(
    description="Convert a mongodb column from string to decimal")

parser.add_argument(
    "--db", type=str, help='mongodb database', required=True)
parser.add_argument(
    "--coll", type=str, help='mongodb collection', required=True)
parser.add_argument(
    '--field', type=str, help='mongodb field to convert', required=True)

args = parser.parse_args()
f, ff = args.field, ''
if '#' in f:
    f, ff = f.split('#')

client = MongoClient()
tbl = client.get_database(args.db).get_collection(args.coll)
print('convert', args.coll + '.' + args.field)
for r in tbl.find(None, [f]):
    dbval = r[f]
    if ff:
        v, need = convert_details(dbval, ff)
    else:
        v, need = convert_value(dbval)
    if need:
        print('.', end='')
        tbl.update_one(
            {'_id': r['_id']},
            {'$set': {f: v}})
    else:
        print(v, end='')
print()
