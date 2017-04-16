$module('mapofanyliteral', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    map = $g.____testlib.basictypes.Mapping($t.struct).overObject(function () {
      var obj = {
      };
      obj['thirdKey'] = $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).overObject(function () {
        var obj = {
        };
        obj['a'] = $t.fastbox('b', $g.____testlib.basictypes.String);
        obj['c'] = $t.fastbox('d', $g.____testlib.basictypes.String);
        return obj;
      }());
      obj['hello'] = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      obj[$t.fastbox(1234, $g.____testlib.basictypes.Integer).String().$wrapped] = $t.fastbox('hello world', $g.____testlib.basictypes.String);
      obj['anotherKey'] = $g.____testlib.basictypes.Slice($g.____testlib.basictypes.Integer).overArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer), $t.fastbox(3, $g.____testlib.basictypes.Integer)]);
      return obj;
    }());
    return $t.fastbox(map.Keys().Length().$wrapped == 4, $g.____testlib.basictypes.Boolean);
  };
});
