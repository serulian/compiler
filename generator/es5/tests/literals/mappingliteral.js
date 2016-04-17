$module('mappingliteral', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.Boolean).overObject(function () {
              var obj = {
              };
              obj['somekey'] = $t.box(true, $g.____testlib.basictypes.Boolean);
              obj['anotherkey'] = $t.box(false, $g.____testlib.basictypes.Boolean);
              return obj;
            }()).then(function ($result0) {
              return $result0.$index($t.box('somekey', $g.____testlib.basictypes.String)).then(function ($result1) {
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
