$module('nochildren', function () {
  var $static = this;
  $static.SimpleFunction = $t.markpromising(function (props, children) {
    var $result;
    var $temp0;
    var $temp1;
    var found;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            found = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
            $current = 1;
            $continue($resolve, $reject);
            return;

          case 1:
            $temp1 = children;
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            value = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              $continue($resolve, $reject);
              return;
            } else {
              $current = 5;
              $continue($resolve, $reject);
              return;
            }
            break;

          case 4:
            found = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 5:
            $resolve(found);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.nochildren.SimpleFunction($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty(), $generator.directempty())).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
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
  });
});
