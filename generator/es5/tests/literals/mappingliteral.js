$module('mappingliteral', function () {
  var $static = this;
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.Boolean).overObject(function () {
              var obj = {
              };
              obj['somekey'] = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
              obj['anotherkey'] = $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
              return obj;
            }()).then(function ($result0) {
              return $result0.$index($t.nominalwrap('somekey', $g.____testlib.basictypes.String)).then(function ($result1) {
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
