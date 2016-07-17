$module('castrejection', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var somevalue;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            somevalue = $t.box('hello', $g.____testlib.basictypes.String);
            $promise.new(function ($resolve) {
              $resolve($t.cast(somevalue, $g.____testlib.basictypes.Integer, false));
            }).then(function ($result0) {
              a = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function ($rejected) {
              b = $rejected;
              $current = 1;
              $continue($resolve, $reject);
              return;
            });
            return;

          case 1:
            $promise.resolve(a == null).then(function ($result0) {
              $result = $t.box($result0 && !(b == null), $g.____testlib.basictypes.Boolean);
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
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
