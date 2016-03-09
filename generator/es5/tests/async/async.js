$module('async', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('0af00909efba2d45fc361ddb814e958e2524c086479947e7b6b5f090c793a4a7', function (a) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.resolve(a);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  });
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $promise.translate($g.async.DoSomethingAsync($t.nominalwrap(3, $g.____testlib.basictypes.Integer))).then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.nominalwrap(3, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $state.current = 1;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $state.resolve($result);
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
