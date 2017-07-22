$module('singlecall', function () {
  var $static = this;
  $static.DoSomething = function () {
    return $t.fastbox(42, $g.________testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    return $t.fastbox($g.singlecall.DoSomething().$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
