$module('mapofanyliteral', function () {
  var $static = this;
  $static.TEST = function () {
    var map;
    map = $g.____testlib.basictypes.Map($t.struct, $t.any).forArrays([$t.fastbox('hello', $g.____testlib.basictypes.String), $t.fastbox(1234, $g.____testlib.basictypes.Integer), $t.fastbox('anotherKey', $g.____testlib.basictypes.String), $t.fastbox('thirdKey', $g.____testlib.basictypes.String)], [$t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox('hello world', $g.____testlib.basictypes.String), $g.____testlib.basictypes.List($g.____testlib.basictypes.Integer).forArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer), $t.fastbox(3, $g.____testlib.basictypes.Integer)]), $g.____testlib.basictypes.Map($g.____testlib.basictypes.String, $g.____testlib.basictypes.String).forArrays([$t.fastbox('a', $g.____testlib.basictypes.String), $t.fastbox('c', $g.____testlib.basictypes.String)], [$t.fastbox('b', $g.____testlib.basictypes.String), $t.fastbox('d', $g.____testlib.basictypes.String)])]);
    return $t.fastbox(map.Mapping().Keys().Length().$wrapped == 4, $g.____testlib.basictypes.Boolean);
  };
});
