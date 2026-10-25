"""October / Halloween puzzle scenes."""
import math, random
from common import *

BAT = 'M0 0 C-8 -10 -22 -14 -40 -10 C-34 -4 -34 2 -38 8 C-30 4 -24 6 -20 12 C-16 6 -10 6 -6 10 C-4 4 -2 2 0 2 C2 2 4 4 6 10 C10 6 16 6 20 12 C24 6 30 4 38 8 C34 2 34 -4 40 -10 C22 -14 8 -10 0 0Z M-4 -2 L-6 -9 L-2 -4 L2 -4 L6 -9 L4 -2Z'


def bat(x, y, s, rot=0, col='#140d1c', op=1):
    return '<path d="%s" fill="%s" opacity="%.2f" transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)"/>' % (BAT, col, op, x, y, rot, s)


def stars(rng, n, ymax, col='#fff8e0'):
    out = []
    for _ in range(n):
        x, y = rng.uniform(0, W), rng.uniform(0, ymax)
        r = rng.choice([1.5, 2, 2.5, 3, 4])
        out.append('<circle cx="%.0f" cy="%.0f" r="%.1f" fill="%s" opacity="%.2f"/>' % (x, y, r, col, rng.uniform(0.4, 1)))
    return ''.join(out)


def pumpkin(x, y, r, rng, face=False, hue=None, glowid='glow'):
    base = hue or rng.choice(['#f07a1a', '#e8691b', '#f28a2a', '#d9601a', '#f39b34'])
    dark = shade(base, -0.35)
    sy = rng.uniform(0.72, 0.85)
    out = ['<g>']
    # shadow
    out.append('<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="#000" opacity="0.25"/>' % (x, y + r*sy*0.92, r*1.1, r*0.18))
    ribs = [(-0.62, 0.55), (0.62, 0.55), (-0.33, 0.62), (0.33, 0.62), (0, 0.6)]
    for i, (ox, w) in enumerate(ribs):
        c = shade(base, -0.12 if i < 2 else (-0.04 if i < 4 else 0.04))
        out.append('<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>' % (
            x + ox*r, y, w*r, r*sy, c, dark, max(1.5, r*0.025)))
    out.append('<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#fff" opacity="0.18"/>' % (x - r*0.25, y - r*sy*0.45, r*0.18, r*sy*0.22))
    # stem
    sx = x + rng.uniform(-0.08, 0.08)*r
    out.append('<path d="M%.1f %.1f q%.1f %.1f %.1f %.1f l%.1f %.1f q%.1f %.1f %.1f %.1f z" fill="#5b6b2a" stroke="#38431a" stroke-width="%.1f"/>' % (
        sx - r*0.09, y - r*sy*0.85, -r*0.02, -r*0.25, r*0.08, -r*0.38, r*0.14, r*0.05, -r*0.08, r*0.15, -r*0.04, r*0.33, max(1, r*0.02)))
    if face:
        g = '#ffd45a'
        fy = y + r*0.05
        e = r*0.2
        out.append('<g filter="url(#%s)" fill="%s">' % (glowid, g))
        for s in (-1, 1):
            out.append('<polygon points="%s"/>' % pts([(x + s*r*0.32 - e*0.6, fy - r*0.12), (x + s*r*0.32 + e*0.6, fy - r*0.12), (x + s*r*0.32, fy - r*0.38)]))
        out.append('<polygon points="%s"/>' % pts([(x - r*0.1, fy + r*0.02), (x + r*0.1, fy + r*0.02), (x, fy - r*0.12)]))
        m = []
        n = 7
        for i in range(n+1):
            t = i/n
            mx = x - r*0.5 + t*r
            m.append((mx, fy + r*0.18 + math.sin(t*math.pi)*r*0.12 + (r*0.08 if i % 2 else 0)))
        for i in range(n, -1, -1):
            t = i/n
            mx = x - r*0.5 + t*r
            m.append((mx, fy + r*0.18 + math.sin(t*math.pi)*r*0.32 - (r*0.06 if i % 2 == 0 else 0)))
        out.append('<polygon points="%s"/>' % pts(m))
        out.append('</g>')
    out.append('</g>')
    return ''.join(out)


def vine(x0, y0, x1, y1, rng, col='#4a5d23'):
    mx, my = (x0+x1)/2, (y0+y1)/2 - rng.uniform(20, 80)
    out = ['<path d="M%.0f %.0f Q%.0f %.0f %.0f %.0f" stroke="%s" stroke-width="7" fill="none" stroke-linecap="round"/>' % (x0, y0, mx, my, x1, y1, col)]
    for _ in range(3):
        t = rng.uniform(0.2, 0.8)
        px = (1-t)**2*x0 + 2*(1-t)*t*mx + t*t*x1
        py = (1-t)**2*y0 + 2*(1-t)*t*my + t*t*y1
        out.append('<path d="M%.0f %.0f c15 -30 45 -20 35 5 c-6 14 -22 10 -18 -2" stroke="%s" stroke-width="4" fill="none"/>' % (px, py, col))
        lx, ly = px + rng.uniform(-10, 10), py
        out.append(leaf_shape(lx, ly - 25, rng.uniform(40, 60), rng.uniform(0, 360), jitter_color('#5f7d2a', rng)))
    return ''.join(out)


def leaf_shape(x, y, s, rot, col):
    d = 'M0 0 C-0.5 -0.3 -0.6 -0.8 0 -1.2 C0.6 -0.8 0.5 -0.3 0 0Z'
    return '<path d="%s" fill="%s" stroke="%s" stroke-width="0.03" transform="translate(%.0f %.0f) rotate(%.0f) scale(%.1f)"/>' % (d, col, shade(col, -0.3), x, y, rot, s)


def dead_tree(x, y, length, angle, width, depth, rng, col):
    out = []

    def br(x, y, l, a, w, d):
        if d == 0 or l < 6:
            return
        x2 = x + math.cos(a)*l
        y2 = y + math.sin(a)*l
        cx = (x+x2)/2 + rng.uniform(-l*0.15, l*0.15)
        cy = (y+y2)/2 + rng.uniform(-l*0.15, l*0.15)
        out.append('<path d="M%.0f %.0f Q%.0f %.0f %.0f %.0f" stroke="%s" stroke-width="%.1f" fill="none" stroke-linecap="round"/>' % (x, y, cx, cy, x2, y2, col, w))
        n = rng.choice([2, 2, 3])
        for i in range(n):
            na = a + rng.uniform(-0.75, 0.75)
            br(x2, y2, l*rng.uniform(0.6, 0.8), na, w*0.68, d-1)
    br(x, y, length, angle, width, depth)
    return ''.join(out)


# ---------------------------------------------------------------- 1
def pumpkin_patch():
    rng = random.Random(101)
    defs = '''<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1">
<stop offset="0" stop-color="#2b1846"/><stop offset="0.45" stop-color="#7a2f5c"/><stop offset="0.75" stop-color="#e2683c"/><stop offset="1" stop-color="#f9b35a"/></linearGradient>
<radialGradient id="moonglow"><stop offset="0" stop-color="#fff3c4" stop-opacity="0.9"/><stop offset="0.3" stop-color="#ffd98a" stop-opacity="0.35"/><stop offset="1" stop-color="#ffd98a" stop-opacity="0"/></radialGradient>
<linearGradient id="field" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#4b3a2a"/><stop offset="1" stop-color="#2a1f17"/></linearGradient>'''
    b = ['<rect width="%d" height="%d" fill="url(#sky)"/>' % (W, H)]
    b.append(stars(rng, 160, 700))
    b.append('<circle cx="2350" cy="430" r="520" fill="url(#moonglow)"/>')
    b.append('<circle cx="2350" cy="430" r="190" fill="#fff1c9"/>')
    for _ in range(9):
        b.append('<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#e9d6a6" opacity="0.5"/>' % (2350 + rng.uniform(-130, 130), 430 + rng.uniform(-130, 130), rng.uniform(12, 34)))
    # clouds
    for i in range(6):
        cx, cy = rng.uniform(0, W), rng.uniform(250, 800)
        b.append('<path d="%s" fill="#c2577a" opacity="0.35" filter="url(#soft)"/>' % blob(cx, cy, rng.uniform(120, 220), rng, 9, 0.3, 2.6, 0.5))
    # hills
    hills = [('#3b2346', 1030, 140), ('#2f1d33', 1120, 110), ('#25172a', 1210, 90)]
    for col, base, amp in hills:
        p = [(0, H)]
        for i in range(0, W+101, 100):
            p.append((i, base - amp*math.sin(i/600 + base) - rng.uniform(0, 30)))
        p.append((W, H))
        b.append('<polygon points="%s" fill="%s"/>' % (pts(p), col))
    # barn silhouette + windmill
    b.append('<g fill="#1d1222"><polygon points="420,1080 420,930 520,860 620,930 620,1080"/><rect x="600" y="960" width="120" height="120"/></g>')
    b.append('<rect x="505" y="960" width="30" height="40" fill="#ffcc66" opacity="0.9"/>')
    b.append(dead_tree(2900, 1150, 170, -math.pi/2, 26, 7, rng, '#1d1222'))
    b.append(dead_tree(180, 1180, 140, -math.pi/2 - 0.1, 22, 6, rng, '#1d1222'))
    # field
    b.append('<path d="M0 1180 Q800 1120 1600 1170 T3200 1150 L3200 2000 L0 2000Z" fill="url(#field)"/>')
    # furrows
    for i in range(14):
        y = 1200 + i*60
        b.append('<path d="M0 %d Q1600 %d 3200 %d" stroke="#1f1712" stroke-width="%d" fill="none" opacity="0.5"/>' % (y, y - 40 + i*3, y, 4 + i))
    # fence
    for i in range(0, W, 140):
        y = 1185 + 12*math.sin(i/300)
        b.append('<rect x="%d" y="%.0f" width="22" height="120" fill="#3a2a22" rx="4"/>' % (i, y - 100))
    b.append('<path d="M0 1120 Q1600 1100 3200 1120" stroke="#4a362a" stroke-width="14" fill="none"/>')
    b.append('<path d="M0 1160 Q1600 1140 3200 1160" stroke="#4a362a" stroke-width="14" fill="none"/>')
    # scarecrow
    sx = 1450
    b.append('<g><rect x="%d" y="880" width="20" height="360" fill="#5a3d26"/><rect x="%d" y="960" width="300" height="18" fill="#5a3d26"/>' % (sx, sx - 140))
    b.append('<path d="M%d 950 L%d 1110 L%d 1110Z" fill="#7b4c8a"/>' % (sx + 10, sx - 70, sx + 90))
    b.append('<path d="M%d 960 L%d 1000 L%d 1000 L%d 960Z" fill="#7b4c8a"/>' % (sx - 130, sx - 130, sx + 150, sx + 150))
    b.append('<circle cx="%d" cy="905" r="50" fill="#e8c88a"/>' % (sx + 10))
    b.append('<path d="M%d 875 L%d 860 L%d 795 L%d 860 L%d 875Z" fill="#3a2a22"/>' % (sx - 70, sx - 40, sx + 10, sx + 60, sx + 90))
    b.append('<circle cx="%d" cy="900" r="7" fill="#222"/><circle cx="%d" cy="900" r="7" fill="#222"/><path d="M%d 925 q20 12 40 0" stroke="#222" stroke-width="4" fill="none"/>' % (sx - 8, sx + 28, sx - 10))
    for k in range(8):
        b.append('<path d="M%d %d l%d %d" stroke="#e7b84a" stroke-width="5"/>' % (sx - 130 + k*3, 970, -20 + k*2, 30))
        b.append('<path d="M%d %d l%d %d" stroke="#e7b84a" stroke-width="5"/>' % (sx + 150 - k*3, 970, 20 - k*2, 30))
    b.append('</g>')
    # crows on fence
    for cx in (900, 980, 2050):
        b.append('<g fill="#120b14" transform="translate(%d 1100)"><ellipse cx="0" cy="-20" rx="26" ry="18"/><circle cx="22" cy="-38" r="13"/><polygon points="32,-40 50,-36 32,-33"/><polygon points="-20,-20 -50,-8 -18,-10"/></g>' % cx)
    # pumpkins, back to front
    items = []
    for _ in range(60):
        y = rng.uniform(1230, 1950)
        t = (y - 1200)/800
        r = 40 + t*170 + rng.uniform(-15, 25)
        x = rng.uniform(-50, W+50)
        items.append((y, x, r))
    items.sort()
    for i, (y, x, r) in enumerate(items):
        if rng.random() < 0.5:
            b.append(vine(x, y, x + rng.uniform(-300, 300), y + rng.uniform(-40, 60), rng))
        face = (r > 120 and rng.random() < 0.45)
        hue = rng.choice([None, None, None, '#f2e6c8', '#e3c148']) if not face else None
        b.append(pumpkin(x, y, r, rng, face=face, hue=hue))
    for i in range(12):
        b.append(bat(rng.uniform(1700, 3000), rng.uniform(200, 800), rng.uniform(0.9, 2.2), rng.uniform(-20, 20)))
    return wrap(defs, ''.join(b))


# ---------------------------------------------------------------- 2
def haunted_house():
    rng = random.Random(202)
    defs = '''<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#0b1030"/><stop offset="0.6" stop-color="#23305e"/><stop offset="1" stop-color="#4a5a8a"/></linearGradient>
<radialGradient id="mg"><stop offset="0" stop-color="#eef3ff" stop-opacity="0.85"/><stop offset="0.35" stop-color="#b8c8f0" stop-opacity="0.3"/><stop offset="1" stop-color="#b8c8f0" stop-opacity="0"/></radialGradient>
<linearGradient id="fog" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#c8d2f0" stop-opacity="0"/><stop offset="1" stop-color="#c8d2f0" stop-opacity="0.28"/></linearGradient>
<linearGradient id="win" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#ffe9a0"/><stop offset="1" stop-color="#ff9c3a"/></linearGradient>'''
    b = ['<rect width="%d" height="%d" fill="url(#sky)"/>' % (W, H), stars(rng, 300, 1100, '#e8eeff')]
    b.append('<circle cx="1600" cy="560" r="700" fill="url(#mg)"/><circle cx="1600" cy="560" r="300" fill="#f4f1e6"/>')
    for _ in range(12):
        b.append('<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#d9d4c2" opacity="0.6"/>' % (1600 + rng.uniform(-220, 220), 560 + rng.uniform(-220, 220), rng.uniform(15, 50)))
    for i in range(7):
        cx, cy = rng.uniform(0, W), rng.uniform(150, 900)
        b.append('<path d="%s" fill="#7d8bb8" opacity="0.35" filter="url(#soft)"/>' % blob(cx, cy, rng.uniform(150, 260), rng, 10, 0.3, 2.8, 0.45))
    # hill
    b.append('<path d="M0 1500 Q500 1250 1100 1180 Q1600 1120 2100 1190 Q2700 1260 3200 1450 L3200 2000 L0 2000Z" fill="#141a33"/>')
    # house
    hx, hy = 1600, 1180
    house = '#0e1226'
    b.append('<g fill="%s">' % house)
    b.append('<rect x="%d" y="%d" width="760" height="420"/>' % (hx - 380, hy - 420))
    b.append('<polygon points="%d,%d %d,%d %d,%d"/>' % (hx - 430, hy - 420, hx, hy - 720, hx + 430, hy - 420))
    b.append('<rect x="%d" y="%d" width="200" height="820"/>' % (hx + 200, hy - 900))
    b.append('<polygon points="%d,%d %d,%d %d,%d"/>' % (hx + 180, hy - 900, hx + 300, hy - 1120, hx + 420, hy - 900))
    b.append('<rect x="%d" y="%d" width="280" height="320"/>' % (hx - 640, hy - 320))
    b.append('<polygon points="%d,%d %d,%d %d,%d"/>' % (hx - 670, hy - 320, hx - 500, hy - 470, hx - 330, hy - 320))
    b.append('<rect x="%d" y="%d" width="50" height="180"/>' % (hx - 250, hy - 760))
    b.append('</g>')
    b.append('<path d="M%d %d l0 -60 l10 0" stroke="%s" stroke-width="6" fill="none"/>' % (hx + 300, hy - 1120, house))
    windows = [(hx - 300, hy - 360), (hx - 160, hy - 360), (hx + 30, hy - 360), (hx - 300, hy - 200), (hx + 30, hy - 200),
               (hx + 260, hy - 820), (hx + 260, hy - 620), (hx + 260, hy - 420), (hx - 580, hy - 250), (hx - 450, hy - 250)]
    for (wx, wy) in windows:
        lit = rng.random() < 0.75
        b.append('<rect x="%d" y="%d" width="80" height="110" rx="4" fill="%s" %s/>' % (wx, wy, 'url(#win)' if lit else '#1f2748', 'filter="url(#glow)"' if lit else ''))
        b.append('<path d="M%d %d v110 M%d %d h80" stroke="%s" stroke-width="6"/>' % (wx + 40, wy, wx, wy + 55, house))
    b.append('<circle cx="%d" cy="%d" r="55" fill="url(#win)" filter="url(#glow)"/><path d="M%d %d h110 M%d %d v110" stroke="%s" stroke-width="6"/>' % (hx, hy - 540, hx - 55, hy - 540, hx, hy - 595, house))
    b.append('<path d="M%d %d v-150 a60 60 0 0 1 120 0 v150Z" fill="#c76b2a" opacity="0.85"/>' % (hx - 100, hy))
    # porch steps
    for i in range(4):
        b.append('<rect x="%d" y="%d" width="%d" height="22" fill="#0b0e1f"/>' % (hx - 130 - i*15, hy + i*22, 180 + i*30))
    # path
    b.append('<path d="M%d %d Q%d %d %d 2000 L%d 2000 Q%d %d %d %dZ" fill="#2a2f4a"/>' % (hx - 70, hy + 88, hx - 120, 1600, hx - 400, hx + 200, hx - 20, 1600, hx + 40, hy + 88))
    # trees
    b.append(dead_tree(500, 1450, 260, -math.pi/2 + 0.08, 46, 8, rng, '#090c1c'))
    b.append(dead_tree(2750, 1420, 240, -math.pi/2 - 0.1, 42, 8, rng, '#090c1c'))
    # graveyard
    for i in range(14):
        gx = rng.uniform(100, 3100)
        if abs(gx - hx) < 520:
            continue
        gy = rng.uniform(1500, 1800)
        s = 0.7 + (gy - 1500)/300*0.8
        stone = jitter_color('#5a6080', rng, 0.06)
        kind = rng.choice(['arch', 'cross', 'rect'])
        if kind == 'arch':
            b.append('<path d="M%.0f %.0f v-%.0f a%.0f %.0f 0 0 1 %.0f 0 v%.0fZ" fill="%s" stroke="#2a2e44" stroke-width="5"/>' % (gx, gy, 110*s, 50*s, 50*s, 100*s, 110*s, stone))
            b.append('<text x="%.0f" y="%.0f" font-family="Lora" font-size="%.0f" fill="#2a2e44" text-anchor="middle">RIP</text>' % (gx + 50*s, gy - 85*s, 30*s))
        elif kind == 'cross':
            b.append('<path d="M%.0f %.0f h%.0f v-%.0f h%.0f v-%.0f h-%.0f v-%.0f h-%.0f v%.0f h-%.0f v%.0f h%.0fZ" fill="%s" stroke="#2a2e44" stroke-width="5"/>' % (
                gx, gy, 30*s, 120*s, 35*s, 30*s, 35*s, 40*s, 30*s, 40*s, 35*s, 30*s, 35*s, stone))
        else:
            b.append('<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="%s" stroke="#2a2e44" stroke-width="5" transform="rotate(%.0f %.0f %.0f)"/>' % (gx, gy - 90*s, 80*s, 90*s, stone, rng.uniform(-8, 8), gx, gy))
        b.append('<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="#0d1124" opacity="0.6"/>' % (gx + 45*s, gy + 4, 80*s, 12*s))
    # iron fence
    for i in range(0, W, 46):
        y = 1880 + 10*math.sin(i/200)
        b.append('<path d="M%d %.0f v-170 l-9 -14 l9 -22 l9 22 l-9 14" fill="#05060f" stroke="#05060f" stroke-width="10"/>' % (i, y + 120))
    b.append('<rect x="0" y="1810" width="3200" height="14" fill="#05060f"/><rect x="0" y="1940" width="3200" height="14" fill="#05060f"/>')
    # fog
    b.append('<rect y="1300" width="3200" height="700" fill="url(#fog)"/>')
    for i in range(10):
        b.append('<path d="%s" fill="#d8def5" opacity="0.18" filter="url(#blur20)"/>' % blob(rng.uniform(0, W), rng.uniform(1450, 1900), rng.uniform(200, 360), rng, 9, 0.2, 3, 0.35))
    for i in range(18):
        b.append(bat(rng.uniform(900, 2300), rng.uniform(200, 900), rng.uniform(0.8, 2.4), rng.uniform(-25, 25), '#05060f'))
    # ghost
    gx, gy = 2350, 1350
    b.append('<g opacity="0.85" filter="url(#glow)"><path d="M%d %d c0 -150 180 -150 180 0 v130 l-30 -25 l-30 25 l-30 -25 l-30 25 l-30 -25 l-30 25Z" fill="#f2f5ff"/><ellipse cx="%d" cy="%d" rx="14" ry="20" fill="#1a1f3a"/><ellipse cx="%d" cy="%d" rx="14" ry="20" fill="#1a1f3a"/><ellipse cx="%d" cy="%d" rx="16" ry="12" fill="#1a1f3a"/></g>' % (gx, gy, gx + 60, gy - 40, gx + 120, gy - 40, gx + 90, gy + 10))
    return wrap(defs, ''.join(b))


# ---------------------------------------------------------------- 3
def maple_path(rng):
    pts_ = []
    lobes = [(0, 1.0), (40, 0.55), (60, 0.85), (95, 0.45), (120, 0.62), (150, 0.25), (180, 0.15)]
    full = lobes + [(360 - a, r) for a, r in reversed(lobes[1:-1])]
    for a, r in full:
        rr = r*(1 + rng.uniform(-0.08, 0.08))
        ang = math.radians(a - 90)
        pts_.append((math.cos(ang)*rr, math.sin(ang)*rr))
    out = []
    for i in range(len(pts_)):
        p, q = pts_[i], pts_[(i+1) % len(pts_)]
        out.append(p)
        mx, my = (p[0]+q[0])/2*0.6, (p[1]+q[1])/2*0.6
        out.append((mx, my))
    return 'M' + ' L'.join('%.3f %.3f' % x for x in out) + 'Z'


def oak_path(rng):
    p = []
    n = 9
    for i in range(n*2+1):
        t = i/(n*2)
        y = -1 + t*1.9
        w = math.sin(t*math.pi)*0.42*(1 + (0.35 if i % 2 else -0.1))
        p.append((w, y))
    left = [(-x, y) for x, y in reversed(p)]
    return smooth_path(p + left)


def autumn_leaves():
    rng = random.Random(303)
    defs = '<linearGradient id="wood" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stop-color="#6b4429"/><stop offset="1" stop-color="#7d5233"/></linearGradient>'
    b = ['<rect width="%d" height="%d" fill="url(#wood)"/>' % (W, H)]
    for i in range(0, H, 250):
        b.append('<rect y="%d" width="%d" height="246" fill="%s"/>' % (i, W, jitter_color('#74492c', rng, 0.05)))
        for k in range(14):
            y = i + rng.uniform(10, 240)
            b.append('<path d="M0 %.0f Q%.0f %.0f 3200 %.0f" stroke="#4e301b" stroke-width="%.0f" opacity="0.35" fill="none"/>' % (y, rng.uniform(800, 2400), y + rng.uniform(-30, 30), y + rng.uniform(-10, 10), rng.uniform(1, 4)))
        b.append('<rect y="%d" width="%d" height="4" fill="#3d2414"/>' % (i + 246, W))
    palette = ['#c0392b', '#d35400', '#e67e22', '#f39c12', '#f1c40f', '#a84a24', '#8e2b1e', '#d4a017', '#b5651d', '#7a8b2c', '#e2552d', '#9c3d54']
    for i in range(420):
        x, y = rng.uniform(-80, W+80), rng.uniform(-80, H+80)
        s = rng.uniform(70, 150)
        rot = rng.uniform(0, 360)
        c = rng.choice(palette)
        c2 = jitter_color(c, rng, 0.12)
        gid = 'lg%d' % i
        defs += '<linearGradient id="%s" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="%s"/><stop offset="1" stop-color="%s"/></linearGradient>' % (gid, shade(c2, 0.15), shade(c, -0.15))
        kind = rng.random()
        if kind < 0.55:
            d = maple_path(rng)
        elif kind < 0.85:
            d = oak_path(rng)
        else:
            d = 'M0 -1 C0.5 -0.6 0.5 0.4 0 0.9 C-0.5 0.4 -0.5 -0.6 0 -1Z'
        veins = '<path d="M0 0.95 L0 -0.85 M0 0.1 L-0.45 -0.35 M0 0.1 L0.45 -0.35 M0 0.45 L-0.3 0.1 M0 0.45 L0.3 0.1" stroke="%s" stroke-width="0.025" fill="none" opacity="0.7"/>' % shade(c, -0.35)
        stem = '<path d="M0 0.9 L0.04 1.35" stroke="%s" stroke-width="0.05"/>' % shade(c, -0.4)
        b.append('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.0f)"><path d="%s" fill="#000" opacity="0.25" transform="translate(0.06 0.08)"/><path d="%s" fill="url(#%s)" stroke="%s" stroke-width="0.02"/>%s%s</g>' % (
            x, y, rot, s, d, d, gid, shade(c, -0.3), veins, stem))
        if i % 35 == 0:
            ax, ay = rng.uniform(0, W), rng.uniform(0, H)
            b.append('<g transform="translate(%.0f %.0f) rotate(%.0f)"><ellipse cx="0" cy="12" rx="34" ry="44" fill="#a0692f" stroke="#5c3a1a" stroke-width="3"/><path d="M-40 -10 Q0 -45 40 -10 Q0 0 -40 -10Z" fill="#6e4a26" stroke="#3e2812" stroke-width="3"/><path d="M0 -30 l4 -20" stroke="#3e2812" stroke-width="6"/></g>' % (ax, ay, rng.uniform(0, 360)))
    return wrap(defs, ''.join(b))


# ---------------------------------------------------------------- 4
def candy_corn(x, y, s, rot):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<path d="M-40 30 Q0 45 40 30 L0 -60Z" fill="#fff4dc" stroke="#5a2a10" stroke-width="3"/>'
            '<path d="M-40 30 Q0 45 40 30 L27 2 Q0 12 -27 2Z" fill="#f6c21a"/>'
            '<path d="M-27 2 Q0 12 27 2 L13 -30 Q0 -26 -13 -30Z" fill="#f07c1c"/>'
            '<path d="M-40 30 Q0 45 40 30 L0 -60Z" fill="none" stroke="#5a2a10" stroke-width="3"/></g>') % (x, y, rot, s)


def wrapped(x, y, s, rot, c):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<path d="M-38 0 L-70 -24 L-64 0 L-70 24Z M38 0 L70 -24 L64 0 L70 24Z" fill="%s" stroke="#2a1530" stroke-width="3"/>'
            '<ellipse rx="42" ry="28" fill="%s" stroke="#2a1530" stroke-width="3"/>'
            '<path d="M-20 -20 L-6 20 M4 -22 L18 18" stroke="#fff" stroke-width="7" opacity="0.6"/></g>') % (x, y, rot, s, shade(c, -0.15), c)


def lollipop(x, y, s, rot, c):
    sp = []
    for i in range(80):
        a = i*0.35; r = i*0.55
        sp.append((math.cos(a)*r, math.sin(a)*r))
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)"><rect x="-4" y="20" width="8" height="90" fill="#f4efe6" stroke="#8a7c6a" stroke-width="2"/>'
            '<circle r="46" fill="%s" stroke="#2a1530" stroke-width="3"/><polyline points="%s" fill="none" stroke="#fff" stroke-width="7" opacity="0.8"/></g>') % (x, y, rot, s, c, pts(sp))


def black_cat(x, y, s, rot):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<path d="M-38 60 C-48 10 -30 -20 -20 -30 L-30 -70 L-8 -45 L8 -45 L30 -70 L20 -30 C30 -20 48 10 38 60Z" fill="#141018"/>'
            '<path d="M38 50 C70 40 80 0 60 -20" stroke="#141018" stroke-width="12" fill="none" stroke-linecap="round"/>'
            '<ellipse cx="-12" cy="-20" rx="8" ry="10" fill="#f6d743"/><ellipse cx="12" cy="-20" rx="8" ry="10" fill="#f6d743"/>'
            '<ellipse cx="-12" cy="-20" rx="2.5" ry="8" fill="#141018"/><ellipse cx="12" cy="-20" rx="2.5" ry="8" fill="#141018"/>'
            '<path d="M-4 -8 L4 -8 L0 -4Z" fill="#e48aa0"/></g>') % (x, y, rot, s)


def ghost(x, y, s, rot):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<path d="M-45 50 L-45 -10 C-45 -75 45 -75 45 -10 L45 50 L30 38 L15 50 L0 38 L-15 50 L-30 38Z" fill="#f8f6ff" stroke="#5b4d78" stroke-width="3"/>'
            '<ellipse cx="-15" cy="-15" rx="7" ry="10" fill="#2a1f3a"/><ellipse cx="15" cy="-15" rx="7" ry="10" fill="#2a1f3a"/>'
            '<ellipse cx="0" cy="8" rx="8" ry="10" fill="#2a1f3a"/><circle cx="-26" cy="0" r="6" fill="#f7a6c0" opacity="0.7"/><circle cx="26" cy="0" r="6" fill="#f7a6c0" opacity="0.7"/></g>') % (x, y, rot, s)


def witch_hat(x, y, s, rot):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<ellipse cx="0" cy="40" rx="70" ry="16" fill="#2b1838" stroke="#120818" stroke-width="3"/>'
            '<path d="M-38 40 L-5 -60 Q10 -70 30 -50 L12 -45 L38 40Z" fill="#3d2350" stroke="#120818" stroke-width="3"/>'
            '<path d="M-33 26 L33 26 L36 38 L-36 38Z" fill="#8bd448"/><rect x="-8" y="24" width="16" height="16" fill="none" stroke="#f6d743" stroke-width="4"/></g>') % (x, y, rot, s)


def cauldron(x, y, s, rot):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<ellipse cx="0" cy="-30" rx="50" ry="14" fill="#7be35a"/>'
            '<circle cx="-15" cy="-50" r="10" fill="#9bf27a" opacity="0.9"/><circle cx="12" cy="-62" r="7" fill="#9bf27a" opacity="0.9"/><circle cx="20" cy="-40" r="6" fill="#9bf27a"/>'
            '<path d="M-55 -30 C-65 40 65 40 55 -30 Z" fill="#26222c" stroke="#0b090e" stroke-width="3"/>'
            '<rect x="-60" y="-36" width="120" height="14" rx="7" fill="#3a3442"/><path d="M-30 30 l-10 20 M30 30 l10 20" stroke="#26222c" stroke-width="10"/></g>') % (x, y, rot, s)


def spider(x, y, s):
    o = ['<g transform="translate(%.0f %.0f) scale(%.2f)"><path d="M0 %d L0 -10" stroke="#d8d0e8" stroke-width="2"/>' % (x, y, s, -int(300/s))]
    for side in (-1, 1):
        for k in range(4):
            o.append('<path d="M0 10 q%d %d %d %d" stroke="#141018" stroke-width="5" fill="none"/>' % (side*25, -20 + k*12, side*45, 0 + k*14))
    o.append('<circle cx="0" cy="10" r="20" fill="#141018"/><circle cx="0" cy="-12" r="12" fill="#141018"/><circle cx="-5" cy="-14" r="3" fill="#ff5050"/><circle cx="5" cy="-14" r="3" fill="#ff5050"/></g>')
    return ''.join(o)


def candy_night():
    rng = random.Random(404)
    defs = '''<radialGradient id="bg" cx="0.5" cy="0.45" r="0.75"><stop offset="0" stop-color="#5b2d82"/><stop offset="1" stop-color="#22103a"/></radialGradient>
<pattern id="dots" width="80" height="80" patternUnits="userSpaceOnUse"><circle cx="20" cy="20" r="4" fill="#7b4aa8" opacity="0.5"/><path d="M60 50 l4 10 l10 2 l-8 7 l2 10 l-8 -5 l-8 5 l2 -10 l-8 -7 l10 -2z" fill="#8a5bb5" opacity="0.35"/></pattern>'''
    b = ['<rect width="%d" height="%d" fill="url(#bg)"/><rect width="%d" height="%d" fill="url(#dots)"/>' % (W, H, W, H)]
    # poisson-ish placement
    placed = []
    cols = ['#e2483d', '#f39c12', '#8e44ad', '#27ae60', '#2e86de', '#ff6fb5', '#f1c40f', '#00b894']
    tries = 0
    while len(placed) < 125 and tries < 8000:
        tries += 1
        x, y = rng.uniform(40, W-40), rng.uniform(40, H-40)
        s = rng.uniform(1.1, 2.3)
        r = 75*s
        if any((x-px)**2 + (y-py)**2 < ((r+pr)*0.82)**2 for px, py, pr in placed):
            continue
        placed.append((x, y, r))
        rot = rng.uniform(-35, 35)
        k = rng.random()
        if k < 0.18: b.append(candy_corn(x, y, s, rot))
        elif k < 0.34: b.append(wrapped(x, y, s, rot*2, rng.choice(cols)))
        elif k < 0.46: b.append(lollipop(x, y, s*0.9, rot, rng.choice(cols)))
        elif k < 0.58: b.append(black_cat(x, y, s*0.9, rot*0.5))
        elif k < 0.69: b.append(ghost(x, y, s, rot*0.5))
        elif k < 0.78: b.append(witch_hat(x, y, s*0.9, rot))
        elif k < 0.86: b.append(cauldron(x, y, s*0.9, rot*0.3))
        elif k < 0.94: b.append(pumpkin(x, y + 10*s, 48*s, rng, face=rng.random() < 0.6))
        else: b.append(bat(x, y, s*1.6, rot, '#120818'))
    for i in range(6):
        b.append(spider(rng.uniform(100, W-100), rng.uniform(150, 600), rng.uniform(1.2, 1.8)))
    return wrap(defs, ''.join(b))


# ---------------------------------------------------------------- 5
def trick_or_treat():
    rng = random.Random(505)
    defs = '''<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#1a1240"/><stop offset="0.55" stop-color="#4b2a6a"/><stop offset="1" stop-color="#c25a3a"/></linearGradient>
<radialGradient id="mg"><stop offset="0" stop-color="#ffb347" stop-opacity="0.8"/><stop offset="0.35" stop-color="#ff8a3a" stop-opacity="0.25"/><stop offset="1" stop-color="#ff8a3a" stop-opacity="0"/></radialGradient>
<linearGradient id="win" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#fff0b0"/><stop offset="1" stop-color="#ffb04a"/></linearGradient>'''
    b = ['<rect width="%d" height="%d" fill="url(#sky)"/>' % (W, H), stars(rng, 200, 800)]
    b.append('<circle cx="700" cy="420" r="600" fill="url(#mg)"/><circle cx="700" cy="420" r="230" fill="#ffc46b"/>')
    for _ in range(10):
        b.append('<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#f0a64f" opacity="0.6"/>' % (700 + rng.uniform(-170, 170), 420 + rng.uniform(-170, 170), rng.uniform(14, 40)))
    # witch on broom across moon
    b.append('<g transform="translate(640 400) rotate(-12)" fill="#120a1c"><rect x="-200" y="20" width="360" height="10" rx="5"/>'
             '<path d="M160 25 L240 -10 L250 25 L240 60Z"/><path d="M-60 25 C-70 -40 -20 -80 0 -80 C20 -60 30 -20 30 25Z"/>'
             '<circle cx="0" cy="-95" r="22"/><path d="M-40 -105 L50 -105 L10 -112 L-20 -190 Z"/><path d="M-60 0 L-120 40 L-40 25Z"/>'
             '<path d="M-20 -60 C-90 -50 -120 -20 -160 -40 C-120 -10 -80 -20 -30 -30Z"/><circle cx="90" cy="5" r="18"/><path d="M80 -10 l-8 -14 l14 8 M98 -10 l8 -14 l-14 8"/></g>')
    # far houses
    for i in range(18):
        x = i*190 - 40 + rng.uniform(-20, 20)
        h = rng.uniform(150, 260)
        col = '#2a1c3e'
        b.append('<rect x="%.0f" y="%.0f" width="170" height="%.0f" fill="%s"/><polygon points="%.0f,%.0f %.0f,%.0f %.0f,%.0f" fill="%s"/>' % (x, 1050 - h, h, col, x - 10, 1050 - h, x + 85, 1050 - h - 80, x + 180, 1050 - h, col))
        for k in range(3):
            if rng.random() < 0.6:
                b.append('<rect x="%.0f" y="%.0f" width="22" height="28" fill="#ffcf73" opacity="0.8"/>' % (x + 20 + k*50, 1050 - h + 40 + rng.choice([0, 60])))
    b.append('<rect y="1050" width="3200" height="950" fill="#1f1530"/>')
    # main houses
    house_cols = ['#5a3d6e', '#3d5a6e', '#6e4a3d', '#4a6e5a', '#6e3d55']
    xs = [-60, 640, 1300, 1980, 2620]
    for i, hx in enumerate(xs):
        w = rng.uniform(560, 640)
        top = rng.uniform(560, 700)
        base = 1500
        col = house_cols[i]
        b.append('<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="%s"/>' % (hx, top, w, base - top, col))
        for yy in range(int(top) + 20, base, 36):
            b.append('<path d="M%.0f %d h%.0f" stroke="%s" stroke-width="3" opacity="0.5"/>' % (hx, yy, w, shade(col, -0.25)))
        b.append('<polygon points="%.0f,%.0f %.0f,%.0f %.0f,%.0f" fill="%s"/>' % (hx - 40, top, hx + w/2, top - rng.uniform(220, 320), hx + w + 40, top, shade(col, -0.45)))
        b.append('<rect x="%.0f" y="%.0f" width="50" height="160" fill="%s"/>' % (hx + w*0.72, top - 230, shade(col, -0.5)))
        for wx in (hx + 70, hx + w - 190):
            b.append('<rect x="%.0f" y="%.0f" width="120" height="150" fill="url(#win)" filter="url(#glow)"/><path d="M%.0f %.0f v150 M%.0f %.0f h120" stroke="%s" stroke-width="8"/>' % (wx, top + 70, wx + 60, top + 70, wx, top + 145, shade(col, -0.5)))
            b.append('<rect x="%.0f" y="%.0f" width="120" height="150" fill="url(#win)" filter="url(#glow)"/><path d="M%.0f %.0f v150 M%.0f %.0f h120" stroke="%s" stroke-width="8"/>' % (wx, top + 330, wx + 60, top + 330, wx, top + 405, shade(col, -0.5)))
        dx = hx + w/2 - 70
        b.append('<rect x="%.0f" y="%.0f" width="140" height="230" fill="%s"/><circle cx="%.0f" cy="%.0f" r="8" fill="#f6d743"/>' % (dx, base - 230, rng.choice(['#8a2c2c', '#2c4a8a', '#2c6a4a', '#6a2c6a']), dx + 115, base - 115))
        b.append('<rect x="%.0f" y="%.0f" width="%.0f" height="24" fill="%s"/>' % (hx - 10, base - 260, w + 20, shade(col, -0.35)))
        for k in range(5):
            px = hx + 20 + k*(w - 40)/4
            b.append('<rect x="%.0f" y="%.0f" width="16" height="250" fill="%s"/>' % (px, base - 250, shade(col, 0.3)))
        # string lights
        lp = []
        for k in range(13):
            t = k/12
            lx = hx + t*w
            ly = top + 20 + math.sin(t*math.pi*4)*14 + 18
            lp.append((lx, ly))
        b.append('<polyline points="%s" stroke="#1a1020" stroke-width="3" fill="none"/>' % pts(lp))
        for k, (lx, ly) in enumerate(lp):
            b.append('<circle cx="%.0f" cy="%.0f" r="10" fill="%s" filter="url(#glow)"/>' % (lx, ly + 10, ['#ff8a2a', '#b45cff', '#7dff6a'][k % 3]))
        # porch pumpkins
        for k in range(rng.choice([2, 3])):
            b.append(pumpkin(hx + rng.uniform(40, w - 40), base + 40, rng.uniform(42, 60), rng, face=True))
    # sidewalk and street
    b.append('<rect y="1540" width="3200" height="160" fill="#4a4258"/>')
    for i in range(0, W, 200):
        b.append('<path d="M%d 1540 v160" stroke="#3a3346" stroke-width="5"/>' % i)
    b.append('<rect y="1700" width="3200" height="300" fill="#2a2433"/>')
    for i in range(0, W, 260):
        b.append('<rect x="%d" y="1840" width="140" height="16" fill="#e8c86a" opacity="0.8"/>' % i)
    # trees with orange leaves
    for tx in (560, 1250, 1930, 2580):
        b.append('<rect x="%d" y="1250" width="40" height="300" fill="#2a1a14"/>' % (tx))
        for k in range(14):
            b.append('<circle cx="%.0f" cy="%.0f" r="%.0f" fill="%s"/>' % (tx + 20 + rng.uniform(-130, 130), 1200 + rng.uniform(-120, 90), rng.uniform(55, 95), rng.choice(['#d35400', '#e67e22', '#c0392b', '#f39c12', '#a84a24'])))
    # fallen leaves on sidewalk + street
    for i in range(160):
        lx, ly = rng.uniform(0, W), rng.uniform(1545, 1995)
        lc = rng.choice(['#d35400', '#e67e22', '#c0392b', '#f39c12', '#a84a24', '#f1c40f'])
        h_ = '<path d="M0 0 C-0.5 -0.3 -0.6 -0.8 0 -1.2 C0.6 -0.8 0.5 -0.3 0 0Z" fill="%s" transform="translate(%.0f %.0f) rotate(%.0f) scale(%.0f)"/>' % (lc, lx, ly, rng.uniform(0, 360), rng.uniform(18, 34))
        b.append(h_)
    b.append('<g transform="translate(2850 1600)"><rect x="-70" y="0" width="140" height="70" rx="10" fill="#5b3a8a"/><text x="0" y="48" font-family="Poppins" font-weight="bold" font-size="34" fill="#ffd45a" text-anchor="middle">CANDY</text><rect x="-60" y="70" width="10" height="60" fill="#3a2a22"/><rect x="50" y="70" width="10" height="60" fill="#3a2a22"/></g>')
    # trick-or-treaters silhouettes
    def kid(x, kind, s=1.0):
        y = 1650
        o = '<g transform="translate(%d %d) scale(%.2f)">' % (x, y, s)
        if kind == 'ghost':
            o += '<path d="M-50 0 L-50 -150 C-50 -230 50 -230 50 -150 L50 0 L30 -15 L10 0 L-10 -15 L-30 0Z" fill="#f4f2fa"/><ellipse cx="-16" cy="-160" rx="9" ry="13" fill="#222"/><ellipse cx="16" cy="-160" rx="9" ry="13" fill="#222"/>'
        elif kind == 'witch':
            o += '<path d="M-45 0 L-30 -140 L30 -140 L45 0Z" fill="#3a1f52"/><circle cx="0" cy="-165" r="32" fill="#9ad66a"/><path d="M-60 -185 L60 -185 L10 -192 L-5 -290Z" fill="#1a0f26"/><rect x="-6" y="-120" width="8" height="150" fill="#5a3d26" transform="rotate(20)"/>'
        elif kind == 'pumpkin':
            o += '<rect x="-25" y="-60" width="18" height="60" fill="#222"/><rect x="7" y="-60" width="18" height="60" fill="#222"/><ellipse cx="0" cy="-110" rx="62" ry="56" fill="#f07a1a"/><circle cx="0" cy="-195" r="30" fill="#f2c9a0"/><rect x="-5" y="-235" width="10" height="16" fill="#3a6a2a"/>'
        else:
            o += '<rect x="-25" y="-70" width="18" height="70" fill="#1a1a1a"/><rect x="7" y="-70" width="18" height="70" fill="#1a1a1a"/><path d="M-45 -70 L-35 -170 L35 -170 L45 -70Z" fill="#1a1a1a"/><path d="M-35 -170 L-90 -60 L-40 -80Z" fill="#a01c2c"/><circle cx="0" cy="-200" r="32" fill="#e8e0f0"/><path d="M-12 -190 l6 10 l6 -10 M6 -190 l6 10 l6 -10" fill="#fff"/>'
        o += '<path d="M45 -100 q40 -10 45 30 q-25 30 -50 0" fill="#e86a1a" stroke="#5a2a10" stroke-width="3"/></g>'
        return o
    b.append(kid(820, 'ghost', 1.05) + kid(980, 'witch', 1.0) + kid(1660, 'pumpkin', 0.9) + kid(1800, 'vamp', 1.05) + kid(2400, 'ghost', 0.85))
    for i in range(10):
        b.append(bat(rng.uniform(1300, 3100), rng.uniform(100, 600), rng.uniform(0.9, 2), rng.uniform(-20, 20), '#120a1c'))
    return wrap(defs, ''.join(b))


SCENES = [
    ('pumpkin-patch', 'Pumpkin Patch at Dusk', pumpkin_patch),
    ('haunted-hill', 'Haunted Hill House', haunted_house),
    ('autumn-leaves', 'Autumn Leaves', autumn_leaves),
    ('candy-night', 'Halloween Candy Night', candy_night),
    ('trick-or-treat', 'Trick-or-Treat Street', trick_or_treat),
]
