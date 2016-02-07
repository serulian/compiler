$module('withas', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var someName;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.nominalwrap(123, $g.____testlib.basictypes.Integer);
            someName = someExpr;
            $state.pushr(someName, 'someName');
            $t.nominalwrap(456, $g.____testlib.basictypes.Integer);
            $state.popr('someName').then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
            $t.nominalwrap(789, $g.____testlib.basictypes.Integer);
            $state.current = -1;
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
