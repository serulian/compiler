$module('asyncchildren', function () {
  var $static = this;
  $static.SimpleFunction = $t.markpromising(function (props, children) {
    var $result;
    var $temp0;
    var $temp1;
    var counter;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            counter = $t.fastbox(0, $g.____testlib.basictypes.Integer);
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
            counter = $t.fastbox(counter.$wrapped + value.$wrapped, $g.____testlib.basictypes.Integer);
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 5:
            $resolve($t.fastbox(counter.$wrapped == 5, $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.DoSomethingAsync = $t.workerwrap('a6e86a18', function (i) {
    return $t.fastbox(i.$wrapped + 1, $g.____testlib.basictypes.Integer);
  });
  $static.GetSomething = $t.markpromising(function (i) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.asyncchildren.DoSomethingAsync(i)).then(function ($result0) {
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
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.asyncchildren.SimpleFunction($g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty(), function () {
              var $result;
              var $current = 0;
              var $continue = function ($yield, $yieldin, $reject, $done) {
                while (true) {
                  switch ($current) {
                    case 0:
                      $promise.maybe($g.asyncchildren.GetSomething($t.fastbox(1, $g.____testlib.basictypes.Integer))).then(function ($result0) {
                        $result = $result0;
                        $current = 1;
                        $continue($yield, $yieldin, $reject, $done);
                        return;
                      }).catch(function (err) {
                        throw err;
                      });
                      return;

                    case 1:
                      $yield($result);
                      $current = 2;
                      return;

                    case 2:
                      $promise.maybe($g.asyncchildren.GetSomething($t.fastbox(2, $g.____testlib.basictypes.Integer))).then(function ($result0) {
                        $result = $result0;
                        $current = 3;
                        $continue($yield, $yieldin, $reject, $done);
                        return;
                      }).catch(function (err) {
                        throw err;
                      });
                      return;

                    case 3:
                      $yield($result);
                      $current = 4;
                      return;

                    default:
                      $done();
                      return;
                  }
                }
              };
              return $generator.new($continue, true);
            }())).then(function ($result0) {
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
