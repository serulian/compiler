$module('loop', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.nominalwrap(1234, $g.____testlib.basictypes.Integer);
            $state.current = 1;
            continue;

          case 1:
            $t.nominalwrap(1357, $g.____testlib.basictypes.Integer);
            $state.current = 1;
            continue;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
