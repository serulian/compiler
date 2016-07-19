$module('resolvereject', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.new(function ($resolve) {
              $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
            }).then(function ($result0) {
              a = $result0;
              b = null;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function ($rejected) {
              b = $rejected;
              a = null;
              $current = 1;
              $continue($resolve, $reject);
              return;
            });
            return;

          case 1:
            $resolve(a);
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
