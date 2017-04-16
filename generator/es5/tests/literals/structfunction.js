$module('structfunction', function () {
  var $static = this;
  $static.buildSomething = function (theMap) {
    return $t.fastbox($t.syncnullcompare(theMap.$index($t.fastbox("Foo", $g.____testlib.basictypes.String)), function () {
      return $t.fastbox(0, $g.____testlib.basictypes.Integer);
    }).$wrapped + $t.syncnullcompare(theMap.$index($t.fastbox("Bar", $g.____testlib.basictypes.String)), function () {
      return $t.fastbox(0, $g.____testlib.basictypes.Integer);
    }).$wrapped, $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    var result;
    result = $g.structfunction.buildSomething($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.Integer).overObject({
      Bar: $t.fastbox(10, $g.____testlib.basictypes.Integer),
      Baz: $t.fastbox(20, $g.____testlib.basictypes.Integer),
      Foo: $t.fastbox(32, $g.____testlib.basictypes.Integer),
    }));
    return $t.fastbox(result.$wrapped == 42, $g.____testlib.basictypes.Boolean);
  };
});
