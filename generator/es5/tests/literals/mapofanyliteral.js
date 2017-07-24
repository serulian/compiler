$module('mapofanyliteral', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    map = $g.________testlib.basictypes.Mapping($t.struct).overObject((function () {
      var obj = {
      };
      obj['thirdKey'] = $g.________testlib.basictypes.Mapping($g.________testlib.basictypes.String).overObject((function () {
        var obj = {
        };
        obj['a'] = $t.fastbox('b', $g.________testlib.basictypes.String);
        obj['c'] = $t.fastbox('d', $g.________testlib.basictypes.String);
        return obj;
      })());
      obj['hello'] = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      obj[$t.fastbox(1234, $g.________testlib.basictypes.Integer).String().$wrapped] = $t.fastbox('hello world', $g.________testlib.basictypes.String);
      obj['anotherKey'] = $g.________testlib.basictypes.Slice($g.________testlib.basictypes.Integer).overArray([$t.fastbox(1, $g.________testlib.basictypes.Integer), $t.fastbox(2, $g.________testlib.basictypes.Integer), $t.fastbox(3, $g.________testlib.basictypes.Integer)]);
      return obj;
    })());
    return $t.fastbox(map.Keys().Length().$wrapped == 4, $g.________testlib.basictypes.Boolean);
  };
});
