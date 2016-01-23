$module('nullcompare', function () {
  var $static = this;
  $static.TEST = function () {
    var someBool;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            someBool = null;
            $state.resolve($t.nullcompare(someBool, true));
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
