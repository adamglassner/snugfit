#!/usr/bin/env python3
"""Add a new puzzle picture to Snugfit.

Usage:
  python tools/add_puzzle.py <image file> "<Puzzle title>" --collection 2026-11 --collection-title "November Collection" [--subtitle "..."]

What it does:
  * resizes/upscales the picture so it is at least 2752 px wide (HD pieces at 1,000 pieces)
  * saves it as assets/puzzles/<id>.jpg and a 480 px thumbnail in assets/puzzles/thumbs/
  * adds it to assets/puzzles/puzzles.json, creating the monthly collection if needed
    (new collections are placed first, so the newest month shows at the top)

Requires: pip install pillow
"""
import argparse, json, os, re, sys
from PIL import Image, ImageFilter

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PUZ = os.path.join(ROOT, 'assets', 'puzzles')


def slug(s):
    return re.sub(r'[^a-z0-9]+', '-', s.lower()).strip('-')


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('image')
    ap.add_argument('title')
    ap.add_argument('--collection', required=True, help='collection id, e.g. 2026-11')
    ap.add_argument('--collection-title', help='e.g. "November Collection" (needed for a new collection)')
    ap.add_argument('--subtitle', default='')
    ap.add_argument('--id', help='puzzle id (default: from title)')
    a = ap.parse_args()

    pid = a.id or slug(a.title)
    im = Image.open(a.image).convert('RGB')
    target_w = 2752
    if im.width < target_w:
        h = round(im.height * target_w / im.width)
        im = im.resize((target_w, h), Image.LANCZOS).filter(ImageFilter.UnsharpMask(radius=2, percent=70, threshold=2))
    elif im.width > 4096:
        im = im.resize((4096, round(im.height * 4096 / im.width)), Image.LANCZOS)
    os.makedirs(os.path.join(PUZ, 'thumbs'), exist_ok=True)
    im.save(os.path.join(PUZ, pid + '.jpg'), quality=92, optimize=True, progressive=True)
    th = im.resize((480, round(im.height * 480 / im.width)), Image.LANCZOS)
    th.save(os.path.join(PUZ, 'thumbs', pid + '.jpg'), quality=88)

    mpath = os.path.join(PUZ, 'puzzles.json')
    m = json.load(open(mpath))
    col = next((c for c in m['collections'] if c['id'] == a.collection), None)
    if col is None:
        if not a.collection_title:
            sys.exit('New collection: please pass --collection-title')
        col = {'id': a.collection, 'title': a.collection_title, 'subtitle': a.subtitle, 'puzzles': []}
        m['collections'].insert(0, col)
    col['puzzles'] = [p for p in col['puzzles'] if p['id'] != pid]
    col['puzzles'].append({'id': pid, 'title': a.title, 'image': pid + '.jpg'})
    json.dump(m, open(mpath, 'w'), indent=2)
    print('Added', pid, 'to', col['title'], '(%dx%d)' % im.size)


if __name__ == '__main__':
    main()
