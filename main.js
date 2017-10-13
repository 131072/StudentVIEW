function newgraph(x, y, z) {
  var bar = new ProgressBar.SemiCircle(x, {
    strokeWidth: 3,
    color: '#FFEA82',
    trailColor: '#eee',
    trailWidth: 1,
    easing: 'easeInOut',
    duration: 1400,
    svgStyle: null,
    text: {
      value: '',
      alignToBottom: false
    },
    from: {color: 'rgb(244, 67, 54)'},
    to: {color: 'rgb(76, 175, 80)'},
    // Set default step function for all animate calls
    step: (state, bar) => {
      bar.path.setAttribute('stroke', state.color);
      var temp = bar.value() * 100
      var value = parseFloat(temp.toFixed(1))
      if (value === 0) {
        bar.setText('');
      } else {
        if (z == undefined) {
          bar.setText(value + "%")
        }
        else {
          bar.setText(value + "% (" + z + ")");
        }
      }

      bar.text.style.color = state.color;
    }
  });
  bar.text.style.fontFamily = '"Raleway", Helvetica, sans-serif';
  bar.text.style.fontSize = '100%';
  x.className = y;
  return bar
}
