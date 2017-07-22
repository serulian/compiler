$module('sliceexpr', function () {
  var $static = this;
  $static.TEST = function () {
    return $g.________testlib.basictypes.Slice($g.________testlib.basictypes.Boolean).overArray([$t.fastbox(false, $g.________testlib.basictypes.Boolean), $t.fastbox(true, $g.________testlib.basictypes.Boolean), $t.fastbox(false, $g.________testlib.basictypes.Boolean)]).$index($t.fastbox(1, $g.________testlib.basictypes.Integer));
  };
});
