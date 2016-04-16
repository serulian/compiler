$module('nullcompare', function () {
  var $static = this;
  $static.TEST = function () {
    var someBool;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            someBool = null;
            $state.resolve($t.nullcompare(someBool, $t.box(true, $g.____testlib.basictypes.Boolean)));
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
