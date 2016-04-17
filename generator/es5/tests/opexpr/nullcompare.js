$module('nullcompare', function () {
  var $static = this;
  $static.TEST = function () {
    var someBool;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      someBool = null;
      $resolve($t.nullcompare(someBool, $t.box(true, $g.____testlib.basictypes.Boolean)));
      return;
    };
    return $promise.new($continue);
  };
});
