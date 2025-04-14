export function fullHeightPageStyleFn(offset: number, height: number) {
  // Based on QPage.style: uses `height` instead of `minHeight`
  return {
    height:
      height === 0
        ? offset === 0
          ? "100vh"
          : `calc(100vh - ${offset}px)`
        : height - offset + "px",
  };
}
