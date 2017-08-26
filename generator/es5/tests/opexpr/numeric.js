$module('numeric', function () {
  var $static = this;
  $static.TEST = function () {
    var result;
    result = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    result = $t.fastbox(result.$wrapped && ((1 + 2) == 3), $g.________testlib.basictypes.Boolean);
    result = $t.fastbox(result.$wrapped && ((2 - 1) == 1), $g.________testlib.basictypes.Boolean);
    result = $t.fastbox(result.$wrapped && ((2 * 1) == 2), $g.________testlib.basictypes.Boolean);
    result = $t.fastbox(result.$wrapped && ($g.________testlib.basictypes.Integer.$div($t.fastbox(5, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer)).$wrapped == 2), $g.________testlib.basictypes.Boolean);
    return result;
  };
});
