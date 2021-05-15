const injectCanva = () => {
  const overwrite = (name) => {
    const OLD = HTMLCanvasElement.prototype[name];
    Object.defineProperty(HTMLCanvasElement.prototype, name, {
      value(...rest) {
        const shift = {
          r: Math.floor(Math.random() * 10) - 5,
          g: Math.floor(Math.random() * 10) - 5,
          b: Math.floor(Math.random() * 10) - 5,
          a: Math.floor(Math.random() * 10) - 5,
        };
        const width = this.width;
        const height = this.height;
        const context = this.getContext('2d');
        const imageData = context.getImageData(0, 0, width, height);
        for (let i = 0; i < height; i += 1) {
          for (let j = 0; j < width; j += 1) {
            const n = i * (width * 4) + j * 4;
            imageData.data[n + 0] = imageData.data[n + 0] + shift.r;
            imageData.data[n + 1] = imageData.data[n + 1] + shift.g;
            imageData.data[n + 2] = imageData.data[n + 2] + shift.b;
            imageData.data[n + 3] = imageData.data[n + 3] + shift.a;
          }
        }
        context.putImageData(imageData, 0, 0);
        return OLD.apply(this, rest);
      },
    });
  };
  overwrite('toBlob');
  overwrite('toDataURL');
};

// eslint-disable-next-line import/prefer-default-export
export { injectCanva };
