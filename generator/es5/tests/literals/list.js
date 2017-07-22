$module('list', function () {
  var $static = this;
  $static.TEST = function () {
    var l;
    l = $g.________testlib.basictypes.Slice($t.struct).overArray([$t.fastbox(1, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer), $t.fastbox(3, $g.________testlib.basictypes.Integer), $t.fastbox(true, $g.________testlib.basictypes.Boolean)]);
    return $t.fastbox(l.Length().$wrapped == 4, $g.________testlib.basictypes.Boolean);
  };
});
