$module('boolean', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            first = true;
            second = false;
            $state.resolve(first && second || first || !second);
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
