$module('structfunction', function () {
  var $static = this;
  $static.buildSomething = function (theMap) {
    return $t.fastbox($t.syncnullcompare(theMap.$index($t.fastbox("Foo", $g.________testlib.basictypes.String)), function () {
      return $t.fastbox(0, $g.________testlib.basictypes.Integer);
    }).$wrapped + $t.syncnullcompare(theMap.$index($t.fastbox("Bar", $g.________testlib.basictypes.String)), function () {
      return $t.fastbox(0, $g.________testlib.basictypes.Integer);
    }).$wrapped, $g.________testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    var result;
    result = $g.structfunction.buildSomething($g.________testlib.basictypes.Mapping($g.________testlib.basictypes.Integer).overObject({
      Bar: $t.fastbox(10, $g.________testlib.basictypes.Integer),
      Baz: $t.fastbox(20, $g.________testlib.basictypes.Integer),
      Foo: $t.fastbox(32, $g.________testlib.basictypes.Integer),
    }));
    return $t.fastbox(result.$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
