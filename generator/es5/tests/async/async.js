$module('async', function () {
  var $static = this;
  $static.DoSomethingAsync = $t.workerwrap('0af00909efba2d45fc361ddb814e958e2524c086479947e7b6b5f090c793a4a7', function (a) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve(a);
      return;
    };
    return $promise.new($continue);
  });
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.async.DoSomethingAsync($t.box(3, $g.____testlib.basictypes.Integer))).then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.box(3, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
