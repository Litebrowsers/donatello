
// sha256 function needs to be here as workers don't share functions from the main thread
async function sha256(uint8Array) {
    const buf = await crypto.subtle.digest("SHA-256", uint8Array);
    return Array.from(new Uint8Array(buf)).map(b => b.toString(16).padStart(2, '0')).join('');
}

class CanvasPredictor {
    constructor(width, height) {
        this.width = width;
        this.height = height;
        this.r = new Uint8Array(width * height).fill(0);
        this.g = new Uint8Array(width * height).fill(0);
        this.b = new Uint8Array(width * height).fill(0);
        this.a = new Uint8Array(width * height).fill(0);
    }

    hexToRGBA(hex) {
        const r = parseInt(hex.substring(0, 2), 16);
        const g = parseInt(hex.substring(2, 4), 16);
        const b = parseInt(hex.substring(4, 6), 16);
        return { r, g, b, a: 255 };
    }

    drawRectangle(r) {
        const col = this.hexToRGBA(r.color);
        for (let y = r.y; y < r.y + r.h; y++) {
            if (y >= 0 && y < this.height) {
                const startX = Math.max(r.x, 0);
                const endX = Math.min(r.x + r.w, this.width);
                if (startX < endX) {
                    const startIndex = y * this.width + startX;
                    const endIndex = y * this.width + endX;
                    this.r.fill(col.r, startIndex, endIndex);
                    this.g.fill(col.g, startIndex, endIndex);
                    this.b.fill(col.b, startIndex, endIndex);
                    this.a.fill(col.a, startIndex, endIndex);
                }
            }
        }
    }

    drawLine(l) {
        const col = this.hexToRGBA(l.color);
        const half = Math.floor(l.thickness / 2);

        if (l.x1 === l.x2) { // Vertical
            const y1 = Math.min(l.y1, l.y2);
            const y2 = Math.max(l.y1, l.y2);
            const startX = Math.max(l.x1 - half, 0);
            const endX = Math.min(l.x1 + half, this.width);

            if(startX < endX) {
                for (let y = y1; y < y2; y++) {
                    if (y >= 0 && y < this.height) {
                        for (let x = startX; x < endX; x++) {
                            const i = y * this.width + x;
                            this.r[i] = col.r;
                            this.g[i] = col.g;
                            this.b[i] = col.b;
                            this.a[i] = col.a;
                        }
                    }
                }
            }
        } else if (l.y1 === l.y2) { // Horizontal
            const x1 = Math.min(l.x1, l.x2);
            const x2 = Math.max(l.x1, l.x2);
            const startY = Math.max(l.y1 - half, 0);
            const endY = Math.min(l.y1 + half, this.height);

            if(startY < endY) {
                for (let y = startY; y < endY; y++) {
                    const startX = Math.max(x1, 0);
                    const endX = Math.min(x2, this.width);
                    if (startX < endX) {
                        const startIndex = y * this.width + startX;
                        const endIndex = y * this.width + endX;
                        this.r.fill(col.r, startIndex, endIndex);
                        this.g.fill(col.g, startIndex, endIndex);
                        this.b.fill(col.b, startIndex, endIndex);
                        this.a.fill(col.a, startIndex, endIndex);
                    }
                }
            }
        }
    }

    drawShapes(taskString) {
        const shapes = taskString.split(';');
        shapes.forEach(shapeStr => {
            if (!shapeStr) return;
            const parts = shapeStr.split(':');
            const type = parts[0];
            const color = parts[1];
            switch (type) {
                case 'R':
                    this.drawRectangle({ color, w: parseInt(parts[2]), h: parseInt(parts[3]), x: parseInt(parts[4]), y: parseInt(parts[5]) });
                    break;
                case 'L':
                    this.drawLine({ color, x1: parseInt(parts[2]), y1: parseInt(parts[3]), x2: parseInt(parts[4]), y2: parseInt(parts[5]), thickness: parseInt(parts[6]) });
                    break;
            }
        });
    }

    async calculateHashes() {
        const [rHash, gHash, bHash, aHash] = await Promise.all([
            sha256(this.r),
            sha256(this.g),
            sha256(this.b),
            sha256(this.a)
        ]);
        return { r: rHash, g: gHash, b: bHash, a: aHash };
    }

    async calculateCombinedHash() {
        const hashes = await this.calculateHashes();
        const combined = hashes.r + hashes.g + hashes.b + hashes.a;
        const combinedHash = await sha256(new TextEncoder().encode(combined));
        return { combinedHash, hashes, channels: { r: this.r, g: this.g, b: this.b, a: this.a } };
    }
}

self.onmessage = async (e) => {
    const { taskString, width, height } = e.data;
    const predictor = new CanvasPredictor(width, height);
    predictor.drawShapes(taskString);
    const result = await predictor.calculateCombinedHash();
    self.postMessage(result);
};
