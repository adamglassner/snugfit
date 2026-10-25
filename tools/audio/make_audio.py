"""Procedurally composes Snugfit's original background music and sound effects.
Run: python3 make_audio.py <out_assets_dir>   (requires numpy, scipy, ffmpeg)"""
import numpy as np, random, subprocess, os, sys
from scipy.signal import fftconvolve, butter, lfilter

SR = 44100
A4 = 440.0
NOTE = {'C': 0, 'C#': 1, 'Db': 1, 'D': 2, 'D#': 3, 'Eb': 3, 'E': 4, 'F': 5, 'F#': 6, 'Gb': 6, 'G': 7, 'G#': 8, 'Ab': 8, 'A': 9, 'A#': 10, 'Bb': 10, 'B': 11}


def midi(name, octave):
    return 12*(octave+1) + NOTE[name]


def freq(m):
    return A4*2**((m-69)/12)


def env_exp(n, attack, tau):
    t = np.arange(n)/SR
    e = np.exp(-t/tau)
    a = int(attack*SR)
    if a > 0:
        e[:a] *= np.linspace(0, 1, a)
    # short release at end
    r = min(n, int(0.03*SR))
    e[-r:] *= np.linspace(1, 0, r)
    return e


def piano(m, dur, vel=0.6):
    f = freq(m)
    n = int((dur + 1.5)*SR)
    t = np.arange(n)/SR
    tau = 1.6*(440/f)**0.35
    tau = min(tau, 3.5)
    sig = np.zeros(n)
    amps = [1, 0.55, 0.32, 0.18, 0.12, 0.07, 0.05]
    for k, a in enumerate(amps, 1):
        fk = f*k*np.sqrt(1 + 0.0004*k*k)
        if fk > 16000:
            break
        sig += a*np.sin(2*np.pi*fk*t + random.random())*np.exp(-t/(tau/(1 + 0.4*(k-1))))
    sig *= env_exp(n, 0.004, 99)
    # damp after note-off
    off = int(dur*SR)
    if off < n:
        sig[off:] *= np.exp(-np.arange(n-off)/SR/0.35)
    return sig*vel*0.35


def musicbox(m, dur, vel=0.6):
    f = freq(m)
    n = int((dur + 2.0)*SR)
    t = np.arange(n)/SR
    sig = np.sin(2*np.pi*f*t)*np.exp(-t/1.4) + 0.35*np.sin(2*np.pi*f*4.07*t)*np.exp(-t/0.25) + 0.15*np.sin(2*np.pi*f*2*t)*np.exp(-t/0.6)
    sig *= env_exp(n, 0.002, 99)
    return sig*vel*0.3


def epiano(m, dur, vel=0.6):
    f = freq(m)
    n = int((dur + 1.2)*SR)
    t = np.arange(n)/SR
    mod = np.sin(2*np.pi*f*t)*1.2*np.exp(-t/0.4)
    sig = np.sin(2*np.pi*f*t + mod)*np.exp(-t/1.6) + 0.2*np.sin(2*np.pi*f*2*t)*np.exp(-t/0.5)
    trem = 1 + 0.08*np.sin(2*np.pi*4.5*t)
    sig *= trem*env_exp(n, 0.003, 99)
    off = int(dur*SR)
    if off < n:
        sig[off:] *= np.exp(-np.arange(n-off)/SR/0.25)
    return sig*vel*0.33


def pad(ms, dur, vel=0.25):
    n = int((dur + 1.0)*SR)
    t = np.arange(n)/SR
    sig = np.zeros(n)
    for m in ms:
        f = freq(m)
        for det in (-0.12, 0, 0.12):
            ff = f*2**(det/12)
            sig += np.sin(2*np.pi*ff*t + random.random()*6) + 0.25*np.sin(2*np.pi*2*ff*t)
    a = int(min(1.2, dur/2)*SR)
    e = np.ones(n)
    e[:a] = np.linspace(0, 1, a)**2
    rel = n - int(dur*SR)
    e[-rel:] = np.linspace(1, 0, rel)
    return sig*e*vel/(len(ms)*3)


def pluck_bass(m, dur, vel=0.5):
    f = freq(m)
    n = int((dur + 0.6)*SR)
    t = np.arange(n)/SR
    sig = (np.sin(2*np.pi*f*t) + 0.3*np.sin(4*np.pi*f*t))*np.exp(-t/0.7)
    return sig*env_exp(n, 0.005, 99)*vel*0.4


def shaker(dur=0.08, vel=0.08):
    n = int(dur*SR)
    s = np.random.randn(n)*np.exp(-np.arange(n)/SR/0.02)
    b, a = butter(2, 5000/(SR/2), 'high')
    return lfilter(b, a, s)*vel


def reverb(x, secs=2.6, wet=0.28, seed=1):
    rng = np.random.default_rng(seed)
    n = int(secs*SR)
    t = np.arange(n)/SR
    irL = rng.standard_normal(n)*np.exp(-t/(secs/5))
    irR = rng.standard_normal(n)*np.exp(-t/(secs/5))
    b, a = butter(1, 6000/(SR/2))
    irL = lfilter(b, a, irL); irR = lfilter(b, a, irR)
    irL /= np.sqrt((irL**2).sum()); irR /= np.sqrt((irR**2).sum())
    L = fftconvolve(x[:, 0], irL)[:len(x)]
    R = fftconvolve(x[:, 1], irR)[:len(x)]
    out = x*(1-wet) + np.stack([L, R], 1)*wet*1.6
    return out


class Track:
    def __init__(self, seconds):
        self.buf = np.zeros((int((seconds + 4)*SR), 2))

    def add(self, sig, start, pan=0.0):
        i = int(start*SR)
        j = min(len(self.buf), i + len(sig))
        s = sig[:j-i]
        l = np.cos((pan + 1)*np.pi/4); r = np.sin((pan + 1)*np.pi/4)
        self.buf[i:j, 0] += s*l
        self.buf[i:j, 1] += s*r

    def finish(self, length, wet=0.28):
        x = reverb(self.buf, wet=wet)
        x = x[:int(length*SR)]
        f = int(3*SR)
        x[-f:] *= np.linspace(1, 0, f)[:, None]
        x[:int(0.5*SR)] *= np.linspace(0, 1, int(0.5*SR))[:, None]
        x /= np.abs(x).max() + 1e-9
        return x*0.85


def scale_notes(root, mode):
    steps = {'major': [0, 2, 4, 5, 7, 9, 11], 'minor': [0, 2, 3, 5, 7, 8, 10], 'dorian': [0, 2, 3, 5, 7, 9, 10]}[mode]
    return [root + 12*o + s for o in range(-2, 4) for s in steps]


def melody_for(chords, scale, rng, beats_per_chord, low, high, density=0.6):
    """returns list of (beat_offset, midi, length_beats)"""
    out = []
    prev = None
    for ci, ch in enumerate(chords):
        b0 = ci*beats_per_chord
        b = 0.0
        while b < beats_per_chord:
            if rng.random() < density or b == 0:
                strong = (b % 2 == 0)
                cands = [n for n in scale if low <= n <= high]
                if strong:
                    cands = [n for n in cands if (n - ch[0]) % 12 in [(c - ch[0]) % 12 for c in ch]]
                if prev is not None:
                    cands.sort(key=lambda n: abs(n - prev) + rng.random()*3)
                    pick = cands[0] if rng.random() < 0.7 else cands[min(len(cands)-1, rng.randint(1, 3))]
                else:
                    pick = rng.choice(cands)
                ln = rng.choice([1, 1, 1.5, 2] if strong else [0.5, 1, 1])
                ln = min(ln, beats_per_chord - b)
                out.append((b0 + b, pick, ln))
                prev = pick
                b += ln
            else:
                b += 0.5
    return out


def song_hearth():
    """Gentle piano, C major, 72bpm"""
    rng = random.Random(11)
    bpm = 72; beat = 60/bpm
    C = lambda n, o: midi(n, o)
    prog_a = [[C('C', 3), C('E', 3), C('G', 3)], [C('A', 2), C('C', 3), C('E', 3)], [C('F', 2), C('A', 2), C('C', 3)], [C('G', 2), C('B', 2), C('D', 3)]]
    prog_b = [[C('F', 2), C('A', 2), C('C', 3)], [C('G', 2), C('B', 2), C('D', 3)], [C('E', 2), C('G', 2), C('B', 2)], [C('A', 2), C('C', 3), C('E', 3)],
              [C('D', 3), C('F', 3), C('A', 3)], [C('G', 2), C('B', 2), C('D', 3)], [C('C', 3), C('E', 3), C('G', 3)], [C('G', 2), C('B', 2), C('D', 3)]]
    sections = [prog_a*2, prog_a*2, prog_b, prog_a*2, prog_b, prog_a*2]
    scale = scale_notes(C('C', 4), 'major')
    total_chords = sum(len(s) for s in sections)
    length = total_chords*4*beat + 4
    tr = Track(length)
    t = 0.0
    motif = None
    for si, sec in enumerate(sections):
        if si == 0:
            mel = []
        elif sec is sections[1] or (si in (1, 3, 5) and motif is not None):
            mel = motif if motif else melody_for(sec, scale, rng, 4, C('E', 4), C('A', 5), 0.55)
        else:
            mel = melody_for(sec, scale, rng, 4, C('E', 4), C('C', 6), 0.6)
        if si == 1:
            motif = mel = melody_for(sec, scale, rng, 4, C('E', 4), C('A', 5), 0.55)
        for ci, ch in enumerate(sec):
            ct = t + ci*4*beat
            # left hand arpeggio
            pattern = [ch[0]-12, ch[1], ch[2], ch[1]+12, ch[2], ch[1]+12, ch[2]+12, ch[1]+12]
            for k, n in enumerate(pattern):
                tr.add(piano(n, beat*0.9, 0.35 + (0.1 if k == 0 else 0) + rng.uniform(-0.04, 0.04)), ct + k*beat/2, -0.25)
            tr.add(pad([ch[0]+12, ch[1]+12, ch[2]+12], 4*beat, 0.10), ct, 0.0)
        for (b, n, ln) in mel:
            tr.add(piano(n, ln*beat*0.95, 0.55 + rng.uniform(-0.05, 0.05)), t + b*beat + rng.uniform(0, 0.012), 0.2)
        t += len(sec)*4*beat
    tr.add(piano(C('C', 2), 4*beat, 0.5) + 0, t)
    tr.add(piano(C('C', 4), 4*beat, 0.4), t)
    tr.add(piano(C('E', 4), 4*beat, 0.4), t + 0.05)
    tr.add(piano(C('G', 4), 4*beat, 0.4), t + 0.1)
    return tr.finish(t + 5, wet=0.3)


def song_lantern():
    """October music box waltz, A minor, 3/4 at 96bpm"""
    rng = random.Random(23)
    bpm = 96; beat = 60/bpm
    m = midi
    prog = [[m('A', 2), m('C', 3), m('E', 3)], [m('A', 2), m('C', 3), m('E', 3)], [m('D', 3), m('F', 3), m('A', 3)], [m('E', 2), m('G#', 2), m('B', 2)],
            [m('F', 2), m('A', 2), m('C', 3)], [m('C', 3), m('E', 3), m('G', 3)], [m('D', 3), m('F', 3), m('A', 3)], [m('E', 2), m('G#', 2), m('B', 2)]]
    prog2 = [[m('F', 2), m('A', 2), m('C', 3)], [m('G', 2), m('B', 2), m('D', 3)], [m('C', 3), m('E', 3), m('G', 3)], [m('A', 2), m('C', 3), m('E', 3)],
             [m('D', 3), m('F', 3), m('A', 3)], [m('B', 2), m('D', 3), m('F', 3)], [m('E', 2), m('G#', 2), m('B', 2)], [m('E', 2), m('G#', 2), m('B', 2)]]
    scale = scale_notes(m('A', 4), 'minor')
    harm = [n for n in scale] + [m('G#', 4), m('G#', 5)]
    sections = [prog, prog, prog2, prog, prog2, prog, prog]
    length = sum(len(s) for s in sections)*3*beat + 5
    tr = Track(length)
    t = 0.0
    mA = melody_for(prog, harm, rng, 3, m('E', 5), m('E', 6), 0.65)
    mB = melody_for(prog2, harm, rng, 3, m('E', 5), m('F', 6), 0.65)
    for si, sec in enumerate(sections):
        mel = [] if si == 0 else (mA if sec is prog else mB)
        for ci, ch in enumerate(sec):
            ct = t + ci*3*beat
            tr.add(pluck_bass(ch[0]-12, beat*1.5, 0.55), ct, -0.1)
            tr.add(musicbox(ch[1]+12, beat, 0.3), ct + beat, -0.35)
            tr.add(musicbox(ch[2]+12, beat, 0.3), ct + beat, -0.35)
            tr.add(musicbox(ch[1]+12, beat, 0.25), ct + 2*beat, -0.35)
            tr.add(musicbox(ch[2]+12, beat, 0.25), ct + 2*beat, -0.35)
            if si >= 2:
                tr.add(pad([ch[0]+12, ch[1]+12, ch[2]+12], 3*beat, 0.07), ct)
        for (b, n, ln) in mel:
            tr.add(musicbox(n, ln*beat, 0.6 + rng.uniform(-0.05, 0.05)), t + b*beat, 0.25)
            if si >= 3 and rng.random() < 0.3:
                tr.add(musicbox(n+12, ln*beat, 0.2), t + b*beat + 0.01, 0.4)
        t += len(sec)*3*beat
    for k, n in enumerate([m('A', 4), m('C', 5), m('E', 5), m('A', 5)]):
        tr.add(musicbox(n, 2, 0.5), t + k*0.18)
    tr.add(pluck_bass(m('A', 1), 3, 0.6), t)
    return tr.finish(t + 5, wet=0.34)


def song_meadow():
    """Warm electric piano, F major, 84bpm with soft shaker"""
    rng = random.Random(37)
    bpm = 84; beat = 60/bpm
    m = midi
    prog = [[m('F', 3), m('A', 3), m('C', 4), m('E', 4)], [m('D', 3), m('F', 3), m('A', 3), m('C', 4)], [m('Bb', 2), m('D', 3), m('F', 3), m('A', 3)], [m('C', 3), m('E', 3), m('G', 3), m('Bb', 3)]]
    prog2 = [[m('G', 3), m('Bb', 3), m('D', 4), m('F', 4)], [m('C', 3), m('E', 3), m('G', 3), m('Bb', 3)], [m('A', 2), m('C', 3), m('E', 3), m('G', 3)], [m('D', 3), m('F', 3), m('A', 3), m('C', 4)],
             [m('Bb', 2), m('D', 3), m('F', 3), m('A', 3)], [m('C', 3), m('E', 3), m('G', 3), m('Bb', 3)], [m('F', 3), m('A', 3), m('C', 4), m('E', 4)], [m('C', 3), m('E', 3), m('G', 3), m('Bb', 3)]]
    scale = scale_notes(m('F', 4), 'major')
    sections = [prog*2, prog*2, prog2, prog*2, prog2, prog*2]
    length = sum(len(s) for s in sections)*4*beat + 5
    tr = Track(length)
    t = 0.0
    mA = melody_for(prog*2, scale, rng, 4, m('A', 4), m('D', 6), 0.5)
    mB = melody_for(prog2, scale, rng, 4, m('C', 5), m('F', 6), 0.55)
    for si, sec in enumerate(sections):
        mel = [] if si == 0 else (mB if sec is prog2 else mA)
        for ci, ch in enumerate(sec):
            ct = t + ci*4*beat
            tr.add(pluck_bass(ch[0]-12, beat*2.5, 0.5), ct, 0)
            tr.add(pluck_bass(ch[0]-12, beat*1, 0.35), ct + 2.5*beat, 0)
            for k in (0, 1.5, 2, 3.5):
                for j, n in enumerate(ch[1:]):
                    tr.add(epiano(n, beat*0.8, 0.28 if k in (0, 2) else 0.18), ct + k*beat + j*0.008, -0.3)
            if si >= 1:
                for k in range(8):
                    tr.add(shaker(vel=0.05 if k % 2 else 0.03), ct + k*beat/2 + rng.uniform(0, 0.01), 0.5)
        for (b, n, ln) in mel:
            tr.add(epiano(n, ln*beat, 0.5 + rng.uniform(-0.05, 0.05)), t + b*beat, 0.2)
        t += len(sec)*4*beat
    for n in prog[0]:
        tr.add(epiano(n, 4, 0.35), t)
    tr.add(pluck_bass(m('F', 1), 4, 0.5), t)
    return tr.finish(t + 5, wet=0.26)


# ---------------------------------------------------------------- SFX
def sfx_tick(f=1800, dur=0.05, noise=0.6, tone=0.5, tau=0.012):
    n = int(dur*SR)
    t = np.arange(n)/SR
    s = noise*np.random.randn(n)*np.exp(-t/(tau*0.6)) + tone*np.sin(2*np.pi*f*t)*np.exp(-t/tau)
    b, a = butter(2, [300/(SR/2), 9000/(SR/2)], 'band')
    return lfilter(b, a, s)


def make_sfx():
    np.random.seed(5)
    out = {}
    out['pick'] = sfx_tick(1500, 0.08, 0.5, 0.6, 0.015)*0.6
    out['drop'] = sfx_tick(700, 0.1, 0.4, 0.8, 0.025)*0.6
    s1 = sfx_tick(2200, 0.06, 0.8, 0.6, 0.01)
    s2 = sfx_tick(1600, 0.07, 0.7, 0.6, 0.012)
    snap = np.zeros(int(0.15*SR)); snap[:len(s1)] += s1; o = int(0.028*SR); snap[o:o+len(s2)] += s2*0.8
    out['snap'] = snap*0.8
    bell = musicbox(midi('E', 6), 0.5, 0.5)[:int(0.9*SR)] + musicbox(midi('B', 6), 0.5, 0.25)[:int(0.9*SR)]
    lock = np.zeros(int(0.9*SR)); lock[:len(snap)] += snap*0.8; lock += bell*0.9
    out['place'] = lock
    comp = np.zeros(int(3.2*SR))
    for k, n in enumerate([midi('C', 5), midi('E', 5), midi('G', 5), midi('C', 6), midi('E', 6), midi('G', 6)]):
        s = musicbox(n, 1.5, 0.6)
        i = int(k*0.11*SR); j = min(len(comp), i+len(s)); comp[i:j] += s[:j-i]
    for k, n in enumerate([midi('C', 4), midi('G', 4), midi('C', 5), midi('E', 5)]):
        s = piano(n, 2.0, 0.4); i = int(0.66*SR); j = min(len(comp), i+len(s)); comp[i:j] += s[:j-i]
    out['complete'] = comp*0.8
    out['click'] = sfx_tick(2600, 0.04, 0.3, 0.7, 0.008)*0.5
    out['hover'] = sfx_tick(3200, 0.025, 0.2, 0.5, 0.005)*0.18
    out['whoosh'] = (lambda n: lfilter(*butter(2, [400/(SR/2), 3000/(SR/2)], 'band'), np.random.randn(n)*np.sin(np.linspace(0, np.pi, n))**2))(int(0.35*SR))*0.25
    return out


def write_ogg(x, path, q=5):
    if x.ndim == 1:
        x = np.stack([x, x], 1)
    x = np.clip(x, -1, 1)
    pcm = (x*32767).astype('<i2').tobytes()
    p = subprocess.run(['ffmpeg', '-y', '-loglevel', 'error', '-f', 's16le', '-ar', str(SR), '-ac', '2', '-i', '-', '-c:a', 'libvorbis', '-q:a', str(q), path], input=pcm)
    assert p.returncode == 0


if __name__ == '__main__':
    out = sys.argv[1]
    os.makedirs(out + '/music', exist_ok=True)
    os.makedirs(out + '/sfx', exist_ok=True)
    for name, fn in [('hearthside', song_hearth), ('lantern-waltz', song_lantern), ('morning-meadow', song_meadow)]:
        random.seed(1); np.random.seed(1)
        x = fn()
        write_ogg(x, '%s/music/%s.ogg' % (out, name), 4)
        print(name, '%.1fs' % (len(x)/SR))
    for k, v in make_sfx().items():
        v = v/ (np.abs(v).max() + 1e-9) * {'hover': 0.35, 'click': 0.6}.get(k, 0.85)
        write_ogg(v, '%s/sfx/%s.ogg' % (out, k), 5)
    print('sfx ok')
